package ziface

// 消息管理抽象层

type IMsgHandle interface {
	// 执行对应的router消息处理方法
	DoMsgHandle(IRequest)
	// 为消息添加具体的处理router
	AddRouter(uint32, IRouter)
}
