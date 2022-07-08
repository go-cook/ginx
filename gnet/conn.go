package gnet

import (
	"fmt"
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
func (c *Conn) SendMsg(msgID uint32, data []byte) error {
	fmt.Println("send msg ", msgID)
	_, err := c.GetTCPConn().Write(data)
	if err != nil {
		fmt.Println("send msg err", err)
		return err
	}
	return nil
}
