package gnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/go-ll/ginx/giface"
)

var defaultHeaderLen uint32 = 8

type Pack struct{}

func NewPack() giface.Pack {
	return &Pack{}
}

func (p *Pack) GetSize() uint32 {
	return defaultHeaderLen
}
func (p *Pack) Pack(msg giface.Message) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	// 写包头Size
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetSize()); err != nil {
		return nil, err
	}
	// 写包头MsgId
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 写包内容
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetContent()); err != nil {
		return nil, err
	}

	// 返回打包内容
	return dataBuff.Bytes(), nil
}
func (p *Pack) UnPack(binaryData []byte) (giface.Message, error) {
	dataBuff := bytes.NewReader(binaryData)
	msg := &Message{}

	// 读头大小
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Size); err != nil {
		return nil, err
	}

	// 读MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.MsgId); err != nil {
		return nil, err
	}

	if msg.Size > 4096 {
		return nil, errors.New("package content is too large")
	}

	return msg, nil
}
