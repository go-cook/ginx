package gnet

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/go-ll/ginx/giface"
)

type Conn struct {
	TCPServer giface.Server
	Conn      *net.TCPConn
	ConnID    uint32
	//消息管理MsgID和对应处理方法的消息管理模块
	MsgHandle giface.MsgHandle
	//告知该链接已经退出/停止的channel
	ctx    context.Context
	cancel context.CancelFunc
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte
	sync.RWMutex
	//链接属性
	property map[string]interface{}
	//保护当前property的锁
	propertyLock sync.Mutex
	IsClosed     bool
}

// NewConn 创建一个业务连接
func NewConn(server giface.Server, conn *net.TCPConn, connID uint32, msgHandle giface.MsgHandle) giface.Conn {
	c := &Conn{
		TCPServer:   server,
		Conn:        conn,
		ConnID:      connID,
		IsClosed:    false,
		MsgHandle:   msgHandle,
		msgBuffChan: make(chan []byte, 1024),
		property:    nil,
	}
	//将新创建的Conn添加到链接管理中
	return c
}

func (c *Conn) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Conn) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Reader exit!]")
	defer c.Stop()

	// 创建拆包解包的对象
	for {
		select {
		case <-c.ctx.Done():
			return
		default:

			//读取客户端的Msg head
			headData := make([]byte, c.TCPServer.Packet().GetSize())
			if _, err := io.ReadFull(c.Conn, headData); err != nil {
				fmt.Println("read msg head error ", err)
				return
			}
			//fmt.Printf("read headData %+v\n", headData)

			//拆包，得到msgID 和 datalen 放在msg中
			msg, err := c.TCPServer.Packet().UnPack(headData)
			if err != nil {
				fmt.Println("unpack error ", err)
				return
			}

			//根据 dataLen 读取 data，放在msg.Data中
			var data []byte
			if msg.GetSize() > 0 {
				data = make([]byte, msg.GetSize())
				if _, err := io.ReadFull(c.Conn, data); err != nil {
					fmt.Println("read msg data error ", err)
					return
				}
			}
			msg.SetContent(data)

			//得到当前客户端请求的Request数据
			req := Request{
				conn: c,
				msg:  msg,
			}
			var WorkerPoolSize int
			WorkerPoolSize = 10

			if WorkerPoolSize > 0 {
				//已经启动工作池机制，将消息交给Worker处理
				c.MsgHandle.SendMsgToTaskQueue(&req)
			} else {
				//从绑定好的消息和对应的处理方法中执行对应的Handle方法
				go c.MsgHandle.DoMsgHandler(&req)
			}
		}
	}
}

func (c *Conn) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TCPServer.CallOnConnStart(c)

	select {
	case <-c.ctx.Done():
		c.finalizer()
		return
	}
}
func (c *Conn) Stop() {
	c.cancel()

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
	c.RLock()
	defer c.RUnlock()
	if c.IsClosed == true {
		return errors.New("connection closed when send msg")
	}

	//将data封包，并且发送
	dp := c.TCPServer.Packet()
	msg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg ID = ", msgID)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	_, err = c.Conn.Write(msg)
	return err
}

func (c *Conn) SendBuffMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	idleTimeout := time.NewTimer(5 * time.Millisecond)
	defer idleTimeout.Stop()

	if c.IsClosed == true {
		return errors.New("connection closed when send buff msg")
	}

	//将data封包，并且发送
	dp := c.TCPServer.Packet()
	msg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg ID = ", msgID)
		return errors.New("Pack error msg ")
	}

	// 发送超时
	select {
	case <-idleTimeout.C:
		return errors.New("send buff msg timeout")
	case c.msgBuffChan <- msg:
		return nil
	}
	//写回客户端
	//c.msgBuffChan <- msg

}

func (c *Conn) SetProperty(key string, value any) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}

	c.property[key] = value
}

func (c *Conn) GetProperty(key string) (any, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}

	return nil, errors.New("no property found")
}

func (c *Conn) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

func (c *Conn) Context() context.Context {
	return c.ctx
}

func (c *Conn) finalizer() {
	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TCPServer.CallOnConnStop(c)

	c.Lock()
	defer c.Unlock()

	//如果当前链接已经关闭
	if c.IsClosed == true {
		return
	}

	fmt.Println("Conn Stop()...ConnID = ", c.ConnID)

	// 关闭socket链接
	_ = c.Conn.Close()

	//将链接从连接管理器中删除
	c.TCPServer.GetConnMgr().Remove(c)

	//关闭该链接全部管道
	close(c.msgBuffChan)
	//设置标志位
	c.IsClosed = true
}
