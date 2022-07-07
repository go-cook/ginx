package gnet

import (
	"net"

	"github.com/go-ll/ginx/giface"
)

type Conn struct {
	TCPServer giface.Server
	Conn      *net.TCPConn
	ConnID    uint32
	IsClosed  bool
}

// NewConn 创建一个业务连接
func NewConn(server giface.Server, conn *net.TCPConn, connID uint32) giface.Conn {
	c := &Conn{
		TCPServer: server,
		Conn:      conn,
		ConnID:    connID,
		IsClosed:  false,
	}
	return c
}

func (c *Conn) Start() {

}
func (c *Conn) Stop() {

}
func (c *Conn) GetTCPConn() *net.TCPConn {
	return c.Conn
}
func (c *Conn) GetConnId() uint32 {
	return c.ConnID
}
func (c *Conn) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
func (c *Conn) SendMsg() {
	// TODO
}
