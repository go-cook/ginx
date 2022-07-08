package giface

type Pack interface {
	GetSize() uint32                           // 包头长度
	Pack(Message) ([]byte, error)              // 打包
	UnPack(binaryData []byte) (Message, error) // 解包
}
