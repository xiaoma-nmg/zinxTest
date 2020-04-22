package ziface

// IRequest接口： 把客户端的连接信息，和请求数据包装在一起

type IRequest interface {
	//得到当前连接
	GetConnection() IConnection
	// 得到请求的消息数据
	GetData() []byte
}
