package ziface

// 将请求的消息，封装到message 中， 定一个抽象的接口

type IMessage interface {
	// 获取消息ID
	GetMsgId() uint32
	// 获取消息的长度
	GetMsgLength() uint32
	// 获取消息内容
	GetData() []byte

	// 设置消息ID
	SetMsgId(uint32)
	// 设置消息长度
	SetMsgLength(uint32)
	// 设置消息内容
	SetData([]byte)
}
