package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 拆包， 封包的单元测试
func TestDataPack(t *testing.T) {
	// 模拟服务器

	// 1创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:9998")
	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}
	// 创建goroutine 承载 负责从客户端处理业务
	go func() {
		// 2 从客户端读取数据，拆包
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}

			go func(conn net.Conn) {
				// 处理客户端的请求
				// 拆包
				// 定义一个拆包的对象
				dp := NewDataPack()
				for {
					// 读取 包头
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error: ", err)
						return
					}

					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack err ", err)
						return
					}

					// 根据 head 中的datalen 再读取 data 的内容
					if msgHead.GetMsgLength() > 0 {
						msg := msgHead.(*Message)
						// 先开辟 Data 需要的空间
						msg.Data = make([]byte, msg.GetMsgLength())

						// 根据dataLen读取
						_, err = io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err: ", err)
							return
						}

						fmt.Println("---> recv MsgID: ", msg.Id, ", datalen = ", msg.DataLen,
							"message is :[", string(msg.Data), "]")
					}
				}

			}(conn)
		}
	}()

	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:9998")
	if err != nil {
		fmt.Println("client dial error: ", err)
		return
	}

	//创建一个 封包对象
	dp := NewDataPack()

	// 模拟粘包过程， 封装两个msg一同发送
	// 第一个msg1
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error: ", err)
		return
	}
	// 第二个msg2
	msg2 := &Message{
		Id:      2,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error: ", err)
		return
	}
	// 沾在一起发送
	sendData1 = append(sendData1, sendData2...)

	_, err = conn.Write(sendData1)
	if err != nil {
		fmt.Println("client send data error: ", err)
		return
	}

	select {}
}
