package znet

import (
	"fmt"
	"io"
	"net"

	"zinx/utils"
	"zinx/ziface"
)

//连接模块
type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 连接的ID
	ConnID uint32
	//当前连接的状态
	isClosed bool
	// 当前连接所绑定的处理业务方法API
	ExitChan chan bool

	// 该连接处理的方法Router
	Router ziface.IRouter
}

//初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		Router:   router,
	}
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running ...")
	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf, 目前最大支持 512 KB
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		cnt, err := c.Conn.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("recv buf err %s \n", err.Error())
			continue
		}

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}
		fmt.Printf("recv data [%s]\n", buf[:cnt])
		// 从路由中，找到注册绑定的Conn对应的router调用
		go func(req ziface.IRequest) {
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.PostHandle(req)
		}(&req)
	}

}

//启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Connection start()... ConnID= ", c.ConnID)
	//启动从当前连接的读数据的业务
	go c.StartReader()

	//TODO 启动从当前连接写数据的业务

}

//停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("conn stop()... ConnID=", c.ConnID)

	if c.isClosed {
		return
	}
	c.isClosed = true

	//关闭socket连接
	_ = c.Conn.Close()
	//关闭channel 回收资源
	close(c.ExitChan)
}

//获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的 TCP的状态 ： IP PORT
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据 给远程的客户端
func (c *Connection) Send([]byte) error {
	return nil
}
