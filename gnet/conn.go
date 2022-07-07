package gnet

import (
	"net"

	"github.com/go-ll/ginx/giface"
)

type Conn struct {
	TCPServer giface.ServerInterface
	Conn      *net.TCPConn
	ConnID    uint32
	IsClosed  bool
}

func NewConn(server giface.ServerInterface, conn *net.TCPConn, connID uint32) *Conn {
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
func (c *Conn) Send() {
	// TODO
}
