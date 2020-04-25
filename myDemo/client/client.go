package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"zinx/znet"
)

func main() {
	fmt.Println("this is a client")
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("Dial server fail")
		return
	}

	for {
		//  把消息打包成指定的TLV格式，发送
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMessage(2, []byte("zinxV0.8 client Test Message")))
		if err != nil {
			fmt.Println("pack error: ", err)
			continue
		}
		_, err = conn.Write(binaryMsg)
		if err != nil {
			fmt.Println("send msg error: ", err)
			continue
		}

		// 接收服务器发来的TLV格式的消息，拆包打印
		//先读取流中的head部分， 得到ID 和 dataLen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error: ", err)
			continue
		}
		msg, err := dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("UnPackage message head error: ", err)
			continue
		}
		// 再根据dataLen 读取
		if msg.GetMsgLength() <= 0 {
			continue
		}
		binaryBody := make([]byte, msg.GetMsgLength())
		if _, err := io.ReadFull(conn, binaryBody); err != nil {
			fmt.Println("read body error: ", err)
			continue
		}
		msg.SetData(binaryBody)
		fmt.Printf("recv service message, messageID:%d, message:%s\n",
			msg.GetMsgId(), string(msg.GetData()))

		// 阻塞cpu
		time.Sleep(time.Second * 2)
	}
}
