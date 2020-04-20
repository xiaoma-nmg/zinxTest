package main

import (
	"zinx/znet"
)

func main() {
	s := znet.NewServer("zinxTest")
	s.Server()
}
