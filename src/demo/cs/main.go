package main

///////////////////////////////////////////
// 将服务端的回包转发到协程中，使其继续执行
// 依据是cmd和seq

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/astaxie/beego"
	"github.com/bbwang/net_frame"
	"github.com/gansidui/gotcp"
)

type SrvCallback struct{}

func (this *SrvCallback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	return true
}

func (this *SrvCallback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	echoPacket := p.(*EchoPacket)
	fmt.Printf("OnMessage:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
	c.AsyncWritePacket(NewEchoPacket(echoPacket.Serialize(), true), time.Second)
	return true
}

func (this *SrvCallback) OnClose(c *gotcp.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}

type CliCallback struct{}

func (this *CliCallback) OnConnect(c *net_frame.Client) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	return true
}

func (this *CliCallback) OnMessage(c *net_frame.Client, p gotcp.Packet) bool {
	echoPacket := p.(*EchoPacket)
	fmt.Printf("OnMessage:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
	c.AsyncWritePacket(NewEchoPacket(echoPacket.Serialize(), true), time.Second)
	return true
}

func (this *CliCallback) OnClose(c *net_frame.Client) {
	fmt.Println("OnClose:", c.RemoteAddr())
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// creates a tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8989")
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// creates a server
	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := gotcp.NewServer(config, &SrvCallback{}, &EchoProtocol{})

	// starts service
	go srv.Start(listener, time.Second)
	fmt.Println("listening:", listener.Addr())

	config2 := &net_frame.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
		Timeout:                time.Second,
	}
	// connect server
	net_frame.NewClient("127.0.0.1:9091", config2, &CliCallback{}, &EchoProtocol{})

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
