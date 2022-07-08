package gnet

import (
	"fmt"
	"net"

	"github.com/go-ll/ginx/giface"
)

type Server struct {
	// 服务名称
	Name string

	// tcpv4
	IPVersion string

	// 服务绑定IP
	IP string

	// 服务绑定端口
	Port int
}

// NewServer 创建一个服务器句柄
func NewServer(name string) giface.Server {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}

func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)

	go func() {
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
			CId++
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("conn read err", err)
						continue
					}
					fmt.Println("conn read count:", cnt)

					conn := NewConn(s, conn, CId)

					if err := conn.SendMsg(conn.GetConnId(), buf); err != nil {
						continue
					}
					//go conn.Start()
				}
			}()
		}
	}()

}

func (s *Server) Stop() {
	fmt.Println("[STOP] Ginx server , name ", s.Name)
}

func (s *Server) Server() {
	s.Start()

	//阻塞,否则主Go退出， listener的go将会退出
	select {}
}

func (s *Server) AddRouter(msgId uint32, router giface.Router) {

}
