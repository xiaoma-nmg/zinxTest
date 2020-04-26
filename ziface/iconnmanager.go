package ziface

// 连接管理模块抽象层

type IConnManager interface {
	Add(conn IConnection)                   //添加连接
	Remove(conn IConnection)                //删除连接
	Get(connID uint32) (IConnection, error) //根据ID获取连接
	Count() int                             //当前总连接数
	ClearConn()                             //清楚并终止所有连接
}
