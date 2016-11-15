package tcplib

import (
	"fmt"
	"net"
)

type Header interface {
	GetLen() uint32
	GetCmd() uint32
}

type TcpRequest struct {
	header  *Header
	body    []byte
	conn    *net.TCPConn
}

func (this *TcpRequest) GetHeader() {
	return header
}

func (this *TcpRequest) Encode() (data []byte, err error) {

}

func (this *TcpRequest) Decode(data []byte) err error {
	
}