package main

import (
	"fmt"

	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

//Handle
func (p *PingRouter) Handle(req ziface.IRequest) {
	fmt.Println("Call PingRouter Handle....")

	// 先读取客户端的数据，再回写 ping...ping...ping...
	fmt.Printf("recv from client: msgID=%d, data=%s\n", req.GetMsgID(), string(req.GetData()))

	msgID := uint32(1)
	if err := req.GetConnection().SendMsg(msgID, []byte("ping...ping...ping...")); err != nil {
		fmt.Printf("send message error, msgID:%d, error:%s\n", msgID, err)
	}
}

type HelloRouter struct {
	znet.BaseRouter
}

//Handle
func (h *HelloRouter) Handle(req ziface.IRequest) {
	fmt.Println("Call HelloRouter Handle....")

	// 先读取客户端的数据，再回写 ping...ping...ping...
	fmt.Printf("recv from client: msgID=%d, data=%s\n", req.GetMsgID(), string(req.GetData()))

	msgID := uint32(2)
	if err := req.GetConnection().SendMsg(msgID, []byte("Hello zinx!!!!")); err != nil {
		fmt.Printf("send message error, msgID:%d, error:%s\n", msgID, err)
	}
}

// 创建连接之后执行的钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnectionBegin is Called...")
	err := conn.SendMsg(202, []byte("DoConnection BEGIN"))
	if err != nil {
		fmt.Println(err)
	}
	// 给当前的连接设置一些属性
	fmt.Println("Set conn Name, Home ...")
	conn.SetProperty("name", "Tom")

}

// 连接断开之前需要执行的函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("DoConnectionList is Called...")
	fmt.Println("Connection ID ", conn.GetConnID(), " is lost...")
	if name, err := conn.GetProperty("name"); err == nil {
		fmt.Println("Name = ", name)
	}

}

func main() {
	// 创建一个server句柄， 使用zinx 的api
	s := znet.NewServer()

	// 注册钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 给当前zinx框架添加一个自定义的router
	s.AddRouter(1, &PingRouter{})
	s.AddRouter(2, &HelloRouter{})

	// 启动server
	s.Server()
}
