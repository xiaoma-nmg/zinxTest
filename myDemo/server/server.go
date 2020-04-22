package main

import (
	"fmt"

	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

//PreHandle
func (p *PingRouter) PreHandle(req ziface.IRequest) {
	fmt.Println("Call Router PreHandle....")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error ", err)
	}
}

//Handle
func (p *PingRouter) Handle(req ziface.IRequest) {
	fmt.Println("Call Router Handle....")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("ping... ping... ping...\n"))
	if err != nil {
		fmt.Println("call back ping error ", err)
	}
}

//PostHandle
func (p *PingRouter) PostHandle(req ziface.IRequest) {
	fmt.Println("Call Router PostHandle....")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back after ping error ", err)
	}
}

func main() {
	// 创建一个server句柄， 使用zinx 的api
	s := znet.NewServer("zinx0.3")

	// 给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	// 启动server
	s.Server()
}
