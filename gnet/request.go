package gnet

import "github.com/go-ll/ginx/giface"

// Request 封装一个请求
type Request struct {
	conn giface.Conn
	msg  giface.Message
}

func (r *Request) GetConn() giface.Conn {
	return r.conn
}
func (r *Request) GetData() []byte {
	return r.msg.GetContent()
}
func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
