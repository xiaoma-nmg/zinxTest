package ziface

type IDataPack interface {
	// 获取包长度的方法
	GetHeadLen() uint32
	// 封包方法
	Pack(IMessage) ([]byte, error)
	// 拆包方法
	UnPack([]byte) (IMessage, error)
}
