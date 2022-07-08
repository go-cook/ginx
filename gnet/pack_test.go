package gnet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestPack(t *testing.T) {
	//创建socket TCP Server
	listener, err := net.Listen("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	//创建服务器gotoutine，负责从客户端goroutine读取粘包的数据，然后进行解析
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept err:", err)
			}

			//处理客户端请求
			go func(conn net.Conn) {
				//创建封包拆包对象dp
				dp := NewPack()
				for {
					//1 先读出流中的head部分
					headData := make([]byte, dp.GetSize())
					_, err := io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止
					if err != nil {
						fmt.Println("read head error")
					}
					//将headData字节流 拆包到msg中
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack err:", err)
						return
					}

					if msgHead.GetSize() > 0 {
						//msg 是有data数据的，需要再次读取data数据
						msg := msgHead.(*Message)
						msg.Content = make([]byte, msg.GetSize())

						//根据dataLen从io中读取字节流
						_, err := io.ReadFull(conn, msg.Content)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							return
						}

						fmt.Println("==> Recv Msg: ID=", msg.MsgId, ", len=", msg.Size, ", data=", string(msg.Content))
					}
				}
			}(conn)
		}
	}()

	//客户端goroutine，负责模拟粘包的数据，然后进行发送
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:7777")
		if err != nil {
			fmt.Println("client dial err:", err)
			return
		}

		//创建一个封包对象 dp
		dp := NewPack()

		//封装一个msg1包
		msg1 := &Message{
			MsgId:   0,
			Size:    5,
			Content: []byte{'h', 'e', 'l', 'l', 'o'},
		}

		sendData1, err := dp.Pack(msg1)
		if err != nil {
			fmt.Println("client pack msg1 err:", err)
			return
		}

		msg2 := &Message{
			MsgId:   1,
			Size:    7,
			Content: []byte{'w', 'o', 'r', 'l', 'd', '!', '!'},
		}
		sendData2, err := dp.Pack(msg2)
		if err != nil {
			fmt.Println("client temp msg2 err:", err)
			return
		}

		//将sendData1，和 sendData2 拼接一起，组成粘包
		sendData1 = append(sendData1, sendData2...)

		//向服务器端写数据
		_, err = conn.Write(sendData1)
		if err != nil {
			return
		}
	}()

	//客户端阻塞
	select {
	case <-time.After(time.Second):
		return
	}
}
