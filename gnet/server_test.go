package gnet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestConnServer(t *testing.T) {
	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		dp := NewPack()
		msg, _ := dp.Pack(NewMsgPackage(1, []byte("client test message")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("client write err: ", err)
			return
		}

		//先读出流中的head部分
		headData := make([]byte, dp.GetSize())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("client read head err: ", err)
			return
		}

		// 将headData字节流 拆包到msg中
		msgHead, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("client unpack head err: ", err)
			return
		}

		if msgHead.GetSize() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			msg := msgHead.(*Message)
			msg.Content = make([]byte, msg.GetSize())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Content)
			if err != nil {
				fmt.Println("client unpack data err")
				return
			}

			fmt.Printf("==> Client receive Msg: ID = %d, len = %d , data = %s\n", msg.MsgId, msg.Size, msg.Content)
		}

		time.Sleep(time.Second)
	}
}

func TestServer(t *testing.T) {
	s := NewServer("Ginx 0.1")
	s.Server()
}
