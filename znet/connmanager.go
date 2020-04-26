package znet

import (
	"errors"
	"fmt"
	"sync"

	"zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理连接集合
	connLock    sync.RWMutex                  //读写锁， 保护连接集合
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (c *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 将连接加入到 管理集 中
	c.connections[conn.GetConnID()] = conn
	fmt.Println("connection ID ", conn.GetConnID(), " add to connManager success! connection count is ", c.Count())
}
func (c *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 将连接从 管理集中删除
	delete(c.connections, conn.GetConnID())
	fmt.Println("connection ID ", conn.GetConnID(), " remove from connManager success! connection count is ", c.Count())
}
func (c *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源，加写锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not exist")
}

func (c *ConnManager) Count() int {
	return len(c.connections)
}

func (c *ConnManager) ClearConn() {
	// 保护共享资源，加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 删除连接，并停止连接的工作
	for connID, conn := range c.connections {
		conn.Stop()
		delete(c.connections, connID)
	}

	fmt.Println("clear all connection success! connection count is ", c.Count())
}
