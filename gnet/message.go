package gnet

type Message struct {
	Size    uint32 // 消息头大小
	MsgId   uint32 // 消息头ID
	Content []byte // 消息内容
}

func NewMsgPackage(Id uint32, data []byte) *Message {
	return &Message{
		Size:    uint32(len(data)),
		MsgId:   Id,
		Content: data,
	}
}

func (m *Message) GetSize() uint32 {
	return m.Size
}
func (m *Message) GetMsgId() uint32 {
	return m.MsgId
}
func (m *Message) GetContent() []byte {
	return m.Content
}

func (m *Message) SetSize(size uint32) {
	m.Size = size
}
func (m *Message) SetMsgId(msgId uint32) {
	m.MsgId = msgId
}
func (m *Message) SetContent(content []byte) {
	m.Content = content
}
