package znet

import (
	"fmt"
	"net"

	"zinx/utils"
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

	//当前的server添加一个router, 是server注册的连接对应的业务处理
	Router ziface.IRouter
}

// 初始化server模块的方法
func NewServer() ziface.IServer {
	return &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name :%s, Listenner at IP:%s, Port:%d\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version :%s, MaxConn:%d, MaxPackageSize:%d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

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
		var cid uint32 = 0
		fmt.Printf("start zinx server %s success, Listinning...\n", s.Name)
		// 阻塞的等待客户端的连接， 处理客户端的读写
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			dealConn := NewConnection(conn, cid, s.Router)
			go dealConn.Start()
			cid++
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

// 添加路由
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
}
