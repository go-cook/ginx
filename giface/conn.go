package giface

import "net"

type Conn interface {
	Start()
	Stop()
	GetTCPConn() *net.TCPConn
	GetConnId() uint32
	RemoteAddr() net.Addr
	SendMsg()
}
