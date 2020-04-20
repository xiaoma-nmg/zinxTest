package znet

import (
	"fmt"
	"io"
	"net"

	"zinx/ziface"
)

//IServer 的接口实现， 定义一个Server的服务器模块
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听IP地址
	IP string
	//服务器监听的端口
	Port int
}

// 初始化server模块的方法
func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      9999,
	}
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP:%s, Port:%d, is starting\n", s.IP, s.Port)
	go func() {
		// 获取一个TCP的ADDR
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}
		// 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, " err ", err)
			return
		}
		fmt.Printf("start zinx server %s success, Listinning...\n", s.Name)
		// 阻塞的等待客户端的连接， 处理客户端的读写
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			// 已经与客户端建立连接，做一个基本的echo业务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err == io.EOF {
						continue
					}
					if err != nil {
						fmt.Println("recv buf err ", err)
						continue
					}
					fmt.Printf("recive client message [%s]\n", buf)

					if _, err = conn.Write(buf[:cnt]); err != nil {
						fmt.Println("echo buf error ", err)
						continue
					}
				}
			}()

		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	//TODO
}

func (s *Server) Server() {
	// 启动server的服务器功能
	s.Start()
	// TODO 做一些启动服务器之后的额外的业务
	//阻塞等待
	select {}
}
