package net_frame

import (
	"net"
)

// common packet interface
type Packet2 interface {
	// 组包，协议格式由包自定义
	Encode() ([]byte, error)
	// 解包，协议格式由包自定义
	Decode(inData []byte) error
}

// common protocol interface
type Protocol2 interface {
	// 根据协议格式分包
	SplitPacket(conn *net.TCPConn) (Packet, error)
}
