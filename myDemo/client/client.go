package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("this is a client")
	addr, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("Dial server fail")
		return
	}

	for {
		_, err := addr.Write([]byte("this is a test from client"))
		if err != nil {
			fmt.Println("write conn err ", err)
			continue
		}

		buf := make([]byte, 512)
		cnt, err := addr.Read(buf)
		if err != nil {
			fmt.Println("read from server error ", err)
			continue
		}
		fmt.Printf("server echo: %s, length=%d\n", buf, cnt)

		time.Sleep(time.Second * 2)
	}
}
