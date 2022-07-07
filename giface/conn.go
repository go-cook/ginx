package giface

import "net"

type ConnectionInterface interface {
	Start()
	Stop()
	GetTCPConn() *net.TCPConn
	GetConnId() uint32
	RemoteAddr() net.Addr
	SendMsg()
}