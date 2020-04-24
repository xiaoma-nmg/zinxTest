package znet

import "zinx/ziface"

type MsgHandler struct {
	Apis map[uint32]ziface.IRouter // 存放每个 MsgID 所对应的消息处理方法
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{Apis: make(map[uint32]ziface.IRouter)}
}

func (m *MsgHandler) DoMsgHandle(req ziface.IRequest) {
	if f, OK := m.Apis[req.GetMsgID()]; OK {
		f.PreHandle(req)
		f.Handle(req)
		f.PostHandle(req)
	}
}

func (m *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, OK := m.Apis[msgID]; OK {
		return
	}

	m.Apis[msgID] = router
}
