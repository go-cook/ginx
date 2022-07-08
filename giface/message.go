package giface

// 封装一个请求消息
type Message interface {
	GetSize() uint32    // 获取消息大小
	GetMsgId() uint32   // 获取消息ID
	GetContent() []byte // 获取消息内容

	SetSize(size uint32)       // 设置消息大小
	SetMsgId(msgId uint32)     // 设置消息ID
	SetContent(content []byte) // 设置消息内容
}
