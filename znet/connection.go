package znet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"zinx/utils"
	"zinx/ziface"
)

//连接模块
type Connection struct {
	Conn      *net.TCPConn      // 当前连接的socket TCP套接字
	ConnID    uint32            // 连接的ID
	isClosed  bool              // 当前连接的状态
	exitChan  chan bool         // 由Reader告知Writer连接退出
	msgChan   chan []byte       // 无缓冲的channel, 用于读写协程之间通信
	MsgHandle ziface.IMsgHandle // 消息的管理msgID, 和对应的处理业务API
	TcpServer ziface.IServer    // 连接属于哪个Server
}

//初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle, server ziface.IServer) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		exitChan:  make(chan bool),
		msgChan:   make(chan []byte),
		MsgHandle: msgHandle,
		TcpServer: server,
	}
	c.TcpServer.GetConnManager().Add(c)
	return c
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running...]")
	defer fmt.Println("connID = ", c.ConnID, "[Reader is exit, remote addr is ]", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 创建一个拆包解包对象
		dp := NewDataPack()
		// 读取客户端的MsgHead
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read message head error: ", err)
			break
		}

		// 拆包，得到MsgID, MsgDataLen
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("Unpack error: ", err)
			break
		}

		// 根据 MsgDataLen 读取MsgBody
		var data []byte
		if msg.GetMsgLength() > 0 {
			data = make([]byte, msg.GetMsgLength())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error: ", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制， 把 req 请求交给 worker pool
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandle.DoMsgHandle(&req)
		}
	}

}

// 专门用于给客户端发送消息
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	// 阻塞等待 Reader 通过channel把数据丢过来，然后发送给客户端
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error: ", err)
				return
			}
		case <-c.exitChan:
			return
		}
	}
}

//启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Connection start()... ConnID= ", c.ConnID)
	//启动从当前连接的读数据的业务
	go c.StartReader()

	// 启动从当前连接写数据的业务
	go c.StartWriter()
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

	// 告知Writer退出
	c.exitChan <- true
	// 将当前连接从连接管理器删除
	c.TcpServer.GetConnManager().Remove(c)
	//关闭channel 回收资源
	close(c.exitChan)
	close(c.msgChan)
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
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection was closed, when send msg")
	}
	// 将data 封包成 message 格式
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgID)
		return errors.New("package error msg")
	}

	// 将数据发送给 Writer
	c.msgChan <- binaryMsg
	return nil
}
