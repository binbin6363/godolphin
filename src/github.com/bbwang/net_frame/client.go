package net_frame

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gansidui/gotcp"
	"net"
	"sync"
	"time"
)

// Error type
var (
	ErrConnClosing          = errors.New("use of closed network connection")
	ErrWriteBlocking        = errors.New("write packet was blocking")
	ErrReadBlocking         = errors.New("read packet was blocking")
	reqSeq           uint32 = 1
)

type Config struct {
	PacketSendChanLimit    uint32        // the limit of packet send channel
	PacketReceiveChanLimit uint32        // the limit of packet receive channel
	Timeout                time.Duration //= (1*time.Second)
}

type NetPacket struct {
	data []byte
}

func (this *NetPacket) Serialize() []byte {
	return this.data
}

func NewPacket(requestData []byte) *NetPacket {
	return &NetPacket{
		data: requestData,
	}
}

// get seq for send request to server
func GetSeq() uint32 {
	reqSeq++
	if 0 == reqSeq {
		reqSeq = 1
	}
	return reqSeq
}

// len, cmd, seq
func GetCmdSeq(data []byte) (cmd, seq uint32) {
	if len(data) < 12 {
		fmt.Println("data len no enough")
		return
	}
	headBuf := bytes.NewBuffer(data[4:12])
	binary.Read(headBuf, binary.BigEndian, &cmd)
	binary.Read(headBuf, binary.BigEndian, &seq)
	return
}

func MakeCmdSeq(cmd, seq uint32) (key uint64) {
	key = uint64(cmd)
	key = (key<<32 | uint64(seq))
	return
}

// 作为client去连接服务端
type Client struct {
	conn              *net.TCPConn                 // the raw connection
	config            *Config                      // Client configuration
	callback          gotcp.ConnCallback           // message callbacks in connection
	protocol          gotcp.Protocol               // customize packet protocol
	exitChan          chan struct{}                // notify all goroutines to shutdown
	waitGroup         *sync.WaitGroup              // wait for all goroutines
	extraData         interface{}                  // to save extra data
	closeOnce         sync.Once                    // close the conn, once, per instance
	closeFlag         int32                        // close flag
	closeChan         chan struct{}                // close chanel
	packetSendChan    chan gotcp.Packet            // packet send chanel
	packetReceiveChan chan gotcp.Packet            // packet receive chanel, unregistered
	recvMap           map[uint64]chan gotcp.Packet // packet receive channel, registered
	timeout           time.Duration                // read timeout
}

// NewClient creates a Client
func NewClient(addr string, config *Config, callback gotcp.ConnCallback, protocol gotcp.Protocol) (*Client, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	client := &Client{
		conn:              conn,
		config:            config,
		callback:          callback,
		protocol:          protocol,
		exitChan:          make(chan struct{}),
		waitGroup:         &sync.WaitGroup{},
		closeChan:         make(chan struct{}),
		packetSendChan:    make(chan gotcp.Packet, config.PacketSendChanLimit),
		packetReceiveChan: make(chan gotcp.Packet, config.PacketReceiveChanLimit),
		recvMap:           make(map[uint64]chan gotcp.Packet),
		timeout:           config.Timeout,
	}
	client.do()
	return client, nil
}

// GetExtraData gets the extra data from the Conn
func (c *Client) GetExtraData() interface{} {
	return c.extraData
}

// PutExtraData puts the extra data with the Conn
func (c *Client) PutExtraData(data interface{}) {
	c.extraData = data
}

// GetRawConn returns the raw net.TCPConn from the Conn
func (c *Client) GetRawConn() *net.TCPConn {
	return c.conn
}

/*func (c *Client) keepAlive() {
    time.AfterFunc()
}
*/
// Start starts service, read response
func (c *Client) readLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.exitChan:
			return

		case <-c.closeChan:
			return

		default:
		}

		p, err := c.protocol.ReadPacket(c.conn)
		if err != nil {
			fmt.Println("read package failed, err:", err.Error())
			return
		}

		cmd, seq := GetCmdSeq(p.Serialize())
		key := MakeCmdSeq(cmd, seq)
		if v, ok := c.recvMap[key]; ok {
			fmt.Println("recv registered response, cmd", cmd, ",seq:", seq)
			v <- p
		} else {
			fmt.Println("recv notice, cmd", cmd, ",seq:", seq)
			c.callback.OnMessage(c, p)
			//c.packetReceiveChan <- p
		}
	}
}

func (c *Client) writeLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.exitChan:
			return

		case <-c.closeChan:
			return

		case p := <-c.packetSendChan:
			if c.IsClosed() {
				return
			}
			if _, err := c.conn.Write(p.Serialize()); err != nil {
				return
			}
		}
	}
}

func (c *Client) GetResult(p []byte) (result []byte, err error) {
	res, err := c.GetResult(NewPacket(p))
	result = res.Serialize()
	return
}

func (c *Client) SendRequest(p []byte) (err error) {
	return c.AsyncWritePacket(NewPacket(p), time.Second)
}

func (c *Client) GetResult(p gotcp.Packet) (result gotcp.Packet, err error) {

	c.packetSendChan <- p
	recvChan := make(chan gotcp.Packet, 1)
	cmd, seq := GetCmdSeq(p.Serialize())

	c.registChan(recvChan, cmd, seq)
	defer func() {
		c.unregistChan(cmd, seq)
		close(recvChan)
	}()

	select {
	case <-c.exitChan:
		err = errors.New("exit chan")
		return

	case <-c.closeChan:
		err = errors.New("close chan")
		return

	case <-time.After(c.timeout):
		err = errors.New("read timeout")
		return

	case result = <-recvChan:
		return
	}

	return
}

func (c *Client) AsyncWritePacket(p gotcp.Packet, timeout time.Duration) (err error) {
	if c.IsClosed() {
		return ErrConnClosing
	}

	defer func() {
		if e := recover(); e != nil {
			err = ErrConnClosing
		}
	}()

	if timeout == 0 {
		select {
		case c.packetSendChan <- p:
			return nil

		// 框架层阻塞写出错
		default:
			return ErrWriteBlocking
		}

	} else {
		select {
		case c.packetSendChan <- p:
			return nil

		// 连接已关闭
		case <-c.closeChan:
			return ErrConnClosing

		// 框架层的超时，代表处理能力不够，大量写请求堆积
		case <-time.After(timeout):
			return ErrWriteBlocking
		}
	}
}

func (this *Client) registChan(recvChan chan gotcp.Packet, cmd, seq uint32) {
	recvKey := MakeCmdSeq(cmd, seq)
	if k, ok := this.recvMap[recvKey]; !ok {
		fmt.Println("regist recvMap cmd:", cmd, ",seq:", seq)
		this.recvMap[recvKey] = recvChan
	} else {
		fmt.Println("recvMap exist cmd:", cmd, ",seq:", seq)
	}
}

func (this *Client) unregistChan(cmd, seq uint32) {
	recvKey := MakeCmdSeq(cmd, seq)
	if k, ok := this.recvMap[recvKey]; ok {
		fmt.Println("unregist chan cmd:", cmd, ",seq:", seq)
		delete(this.recvMap, recvKey)
	} else {
		fmt.Println("unregist chan, key not exist, cmd:", cmd, ",seq:", seq)
	}
}

// Do it, 链式调用
func (c *Client) do() *Client {
	if !c.callback.OnConnect(c) {
		return
	}

	asyncDo(c.readLoop, c.waitGroup)
	asyncDo(c.writeLoop, c.waitGroup)
	return c
}

func asyncDo(fn func(), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		fn()
		wg.Done()
	}()
}

// Stop stops service
func (s *Client) Stop() {
	close(s.exitChan)
	s.waitGroup.Wait()
}

// Close closes the connection
func (c *Client) Close() {
	c.closeOnce.Do(func() {
		atomic.StoreInt32(&c.closeFlag, 1)
		close(c.closeChan)
		close(c.packetSendChan)
		close(c.packetReceiveChan)
		c.conn.Close()
		c.callback.OnClose(c)
	})
}

// IsClosed indicates whether or not the connection is closed
func (c *Client) IsClosed() bool {
	return atomic.LoadInt32(&c.closeFlag) == 1
}
