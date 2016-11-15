package gotcp

import (
	"github.com/bbwang/logs"
)

/**
*   包括自动重连机制，心跳机制
**/

// connection, connect to server
type ServerSession struct {
}

// 连接到指定的服务，可指定协议
// addr: host:port, eg: 127.0.0.1:9000
// protocol: tcp/http
func (this *ServerSession) Connect(addr string, protocol int) (session *ServerSession, err error) {

}
