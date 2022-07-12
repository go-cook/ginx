package gnet

import (
	"fmt"
	"net"

	"github.com/go-ll/ginx/giface"
)

var ginxLogo = `                                        
  ____ _
 / ___(_)_ __ __  __
| |  _| | '_ \\ \/ /
| |_| | | | | |>  <
 \____|_|_| |_/_/\_\
                                        `
var topLine = `┌──────────────────────────────────────────────────────┐`
var borderLine = `│`
var bottomLine = `└──────────────────────────────────────────────────────┘`

type Server struct {
	//服务器的名称
	Name string
	//tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port int
	//当前Server的消息管理模块，用来绑定MsgID和对应的处理方法
	msgHandler giface.MsgHandle
	//当前Server的链接管理器
	ConnMgr giface.ConnManager
	//该Server的连接创建时Hook函数
	OnConnStart func(conn giface.Conn)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn giface.Conn)

	packet giface.Pack
}

// NewServer 创建一个服务器句柄
func NewServer(opts ...Option) giface.Server {
	printLogo()
	s := &Server{
		Name:       "ginx",
		IPVersion:  "tcp4",
		IP:         "0.0.0.0",
		Port:       8999,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
		packet:     NewPack(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)

	go func() {
		//0 启动worker工作池机制
		s.msgHandler.StartWorkerPool()
		// 1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resole tcp addr err:", err)
			return

		}
		// 2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			panic(err)
		}

		//已经监听成功
		fmt.Println("start ginx server  ", s.Name, " succ, now listenning...")

		var CId uint32
		CId = 0

		// 3 启动server网络连接业务
		for {
			// 3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			fmt.Println("Get conn remote addr = ", conn.RemoteAddr().String())

			//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Count() >= 1024 {
				conn.Close()
				continue
			}

			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConn(s, conn, CId, s.msgHandler)
			CId++

			//3.4 启动当前链接的处理业务
			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	fmt.Println("[STOP] Ginx server , name ", s.Name)
	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	//阻塞,否则主Go退出， listener的go将会退出
	select {}
}

//AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgId uint32, router giface.Router) {
	s.msgHandler.AddRouter(msgId, router)
}

//GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() giface.ConnManager {
	return s.ConnMgr
}

//SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(giface.Conn)) {
	s.OnConnStart = hookFunc
}

//SetOnConnStop 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(giface.Conn)) {
	s.OnConnStop = hookFunc
}

//CallOnConnStart 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn giface.Conn) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

//CallOnConnStop 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn giface.Conn) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}

func (s *Server) Packet() giface.Pack {
	return s.packet
}

func printLogo() {
	fmt.Println(ginxLogo)
	fmt.Println(topLine)
	fmt.Println(fmt.Sprintf("%s [Github] https://github.com/go-ll/ginx                    %s", borderLine, borderLine))
	fmt.Println(bottomLine)
	fmt.Printf("[Ginx] Version: %f, MaxConn: %d, MaxPacketSize: %d\n",
		1.0,
		1024,
		1024)
}

func init() {
}
