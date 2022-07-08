package giface

type Request interface {
	GetConn() Conn    // 请求连接
	GetData() []byte  // 请求内容
	GetMsgId() uint32 // 请求ID
}
