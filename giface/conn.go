package giface

import (
	"context"
	"net"
)

type Conn interface {
	Start()
	Stop()
	//返回ctx，用于用户自定义的go程获取连接退出状态

	Context() context.Context

	GetTCPConn() *net.TCPConn
	GetConnId() uint32
	RemoteAddr() net.Addr

	SendMsg(msgID uint32, data []byte) error
	SendBuffMsg(msgId uint32, data []byte) error

	// SetProperty 设置链接属性
	SetProperty(key string, value any)
	// GetProperty 获取链接属性
	GetProperty(key string) (any, error)
	RemoveProperty(key string)
}
