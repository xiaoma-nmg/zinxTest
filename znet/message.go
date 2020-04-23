package znet

type Message struct {
	Id      uint32 //消息的ID
	DataLen uint32 // 消息长度
	Data    []byte //消息内容
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// 获取消息ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

// 获取消息的长度
func (m *Message) GetMsgLength() uint32 {
	return m.DataLen
}

// 获取消息内容
func (m *Message) GetData() []byte {
	return m.Data
}

// 设置消息ID
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

// 设置消息长度
func (m *Message) SetMsgLength(length uint32) {
	m.DataLen = length
}

// 设置消息内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}
