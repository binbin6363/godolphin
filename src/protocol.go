package gotcp

import (
	"net"
)

// common packet interface
type Packet interface {
	Serialize() []byte
}

// common protocol interface
type Protocol interface {
	ReadPacket(conn *net.TCPConn) (Packet, error)
}



