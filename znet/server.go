package znet

import (
	"fmt"
	"net"

	"zinx/utils"
	"zinx/ziface"
)

//IServer 的接口实现， 定义一个Server的服务器模块
type Server struct {
	Name        string                        //服务器的名称
	IPVersion   string                        //服务器绑定的ip版本
	IP          string                        //服务器监听IP地址
	Port        int                           //服务器监听的端口
	MsgHandle   ziface.IMsgHandle             //绑定msgID 和对应的处理业务
	ConnMgr     ziface.IConnManager           //服务器的连接管理器
	OnConnStart func(conn ziface.IConnection) // server创建后自动调用的hook函数
	OnConnStop  func(conn ziface.IConnection) // server销毁连接之前
}

// 初始化server模块的方法
func NewServer() ziface.IServer {
	return &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		MsgHandle: NewMsgHandler(),
		ConnMgr:   NewConnManager(),
	}
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name :%s, Listenner at IP:%s, Port:%d\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version :%s, MaxConn:%d, MaxPackageSize:%d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	go func() {
		// 开启消息队列及 worker 工作池
		s.MsgHandle.StartWorkerPool()

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

			// 判断连接个数是否超过最大值
			if s.ConnMgr.Count() >= utils.GlobalObject.MaxConn {
				fmt.Println("该服务器当前连接数已达到最大值")
				conn.Close()
				continue
			}

			dealConn := NewConnection(conn, cid, s.MsgHandle, s)
			go dealConn.Start()
			cid++
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	//TODO
	fmt.Println("[STOP] zinx server ", s.Name, " stop ")
	s.ConnMgr.ClearConn()
}

func (s *Server) Server() {
	// 启动server的服务器功能
	s.Start()
	// TODO 做一些启动服务器之后的额外的业务
	//阻塞等待
	select {}
}

// 添加路由
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandle.AddRouter(msgID, router)
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnMgr
}

// 注册OnConnStart 钩子函数
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 注册OnConnStop 钩子函数
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用OnConnStart 钩子函数
func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("Call OnConnStart()... ")
		s.OnConnStart(connection)
	}
}

// 调用OnConnStart 钩子函数
func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("Call OnConnStop()... ")
		s.OnConnStop(connection)
	}
}
