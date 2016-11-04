package gotcp

import (
	"errors"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"github.com/astaxie/beego/logs"
)

// Error type
var (
	ErrConnClosing   = errors.New("use of closed network connection")
	ErrWriteBlocking = errors.New("write packet was blocking")
	ErrReadBlocking  = errors.New("read packet was blocking")
)

// Conn exposes a set of callbacks for the various events that occur on a connection
type Conn struct {
	srv               *Server
	conn              *net.TCPConn  // the raw connection
	extraData         interface{}   // to save extra data
	closeOnce         sync.Once     // close the conn, once, per instance
	closeFlag         int32         // close flag
	closeChan         chan struct{} // close chanel
	packetSendChan    chan Packet   // packet send chanel
	packetReceiveChan chan Packet   // packeet receive chanel
}

// ConnCallback is an interface of methods that are used as callbacks on a connection
type ConnCallback interface {
	// OnConnect is called when the connection was accepted,
	// If the return value of false is closed
	OnConnect(*Conn) bool

	// OnMessage is called when the connection receives a packet,
	// If the return value of false is closed
	OnMessage(*Conn, Packet) bool

	// OnClose is called when the connection closed
	OnClose(*Conn)
}

// newConn returns a wrapper of raw conn
func newConn(conn *net.TCPConn, srv *Server) *Conn {
	return &Conn{
		srv:               srv,
		conn:              conn,
		closeChan:         make(chan struct{}),
		packetSendChan:    make(chan Packet, srv.config.PacketSendChanLimit),
		packetReceiveChan: make(chan Packet, srv.config.PacketReceiveChanLimit),
	}
}

// GetExtraData gets the extra data from the Conn
func (c *Conn) GetExtraData() interface{} {
	return c.extraData
}

// PutExtraData puts the extra data with the Conn
func (c *Conn) PutExtraData(data interface{}) {
	c.extraData = data
}

// GetRawConn returns the raw net.TCPConn from the Conn
func (c *Conn) GetRawConn() *net.TCPConn {
	return c.conn
}

// Close closes the connection
func (c *Conn) Close() {
	c.closeOnce.Do(func() {
		atomic.StoreInt32(&c.closeFlag, 1)
		close(c.closeChan) // 关闭连接
		close(c.packetSendChan)
		close(c.packetReceiveChan)
		c.conn.Close()
		c.srv.callback.OnClose(c)
	})
}

// IsClosed indicates whether or not the connection is closed
func (c *Conn) IsClosed() bool {
	return atomic.LoadInt32(&c.closeFlag) == 1
}

// AsyncWritePacket async writes a packet, this method will never block
func (c *Conn) AsyncWritePacket(p Packet, timeout time.Duration) (err error) {
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

		// 收到连接关闭的信号
		case <-c.closeChan:
			return ErrConnClosing

		// 框架层的超时，代表处理能力不够，大量写请求堆积
		case <-time.After(timeout):
			return ErrWriteBlocking
		}
	}
}

// Do it, 链式调用
func (c *Conn) Do() *Conn {
	if !c.srv.callback.OnConnect(c) {
		return
	}

	asyncDo(c.handleLoop, c.srv.waitGroup)
	asyncDo(c.readLoop, c.srv.waitGroup)
	asyncDo(c.writeLoop, c.srv.waitGroup)
	return c
}

// 从网络读取数据并分包传到接收队列
func (c *Conn) readLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.srv.exitChan:
			return

		// 收到连接关闭的信号
		case <-c.closeChan:
			return

		default:
		}

		p, err := c.srv.protocol.SplitPacket(c.conn)
		if err != nil {
			// TODO: 包错误，考虑到安全因素，需要断开连接，断开之前是否要等待写数据包完成呢？
			select{
				case <- time.After(1*time.Second)
				Close()
				break
			}
			return
		}

		c.packetReceiveChan <- p
	}
}

// 循环写队列中的数据到网络
func (c *Conn) writeLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.srv.exitChan:
			return

		// 收到连接关闭的信号
		case <-c.closeChan:
			return

		case p := <-c.packetSendChan:
			if c.IsClosed() {
				return
			}
			if data, err := p.Encode(); err != nil {
				if _, err = c.conn.Write(data); err != nil {
					log.Warn("write data failed, disconnect net connection.")
					Close()
					// 写出错就不用等待读取完成
					return
				}
			} else {
				log.Warn("encode data failed.")
			}
		}
	}
}

// 取接收队列的包，处理业务逻辑
func (c *Conn) handleLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.srv.exitChan:
			return

		// 收到连接关闭的信号
		case <-c.closeChan:
			return

		case p := <-c.packetReceiveChan:
			if c.IsClosed() {
				return
			}
			if !c.srv.callback.OnMessage(c, p) {
				return
			}
		}
	}
}

func asyncDo(fn func(), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		fn()
		wg.Done()
	}()
}
