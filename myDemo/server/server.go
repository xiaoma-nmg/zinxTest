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
	fmt.Println("Call Router Handle....")

	// 先读取客户端的数据，再回写 ping...ping...ping...
	fmt.Printf("recv from client: msgID=%d, data=%s\n", req.GetMsgID(), string(req.GetData()))

	msgID := uint32(1)
	if err := req.GetConnection().SendMsg(msgID, []byte("ping...ping...ping...")); err != nil {
		fmt.Printf("send message error, msgID:%d, error:%s\n", msgID, err)
	}
}

func main() {
	// 创建一个server句柄， 使用zinx 的api
	s := znet.NewServer()

	// 给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	// 启动server
	s.Server()
}
