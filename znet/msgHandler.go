package znet

import (
	"fmt"

	"zinx/utils"
	"zinx/ziface"
)

// 消息处理模块的实现
type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter // 存放每个 MsgID 所对应的消息处理方法
	TaskQueue      []chan ziface.IRequest    // 负责worker 取任务的消息队列
	WorkerPoolSize uint32                    // worker 工作池的数量
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, // 从全局配置中获取
	}
}

// 执行对应的msgID的 Handle 函数
func (m *MsgHandler) DoMsgHandle(req ziface.IRequest) {
	if f, OK := m.Apis[req.GetMsgID()]; OK {
		f.PreHandle(req)
		f.Handle(req)
		f.PostHandle(req)
	}
}

// 为消息添加 具体的处理逻辑
func (m *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, OK := m.Apis[msgID]; OK {
		return
	}

	m.Apis[msgID] = router
}

// 启动一个worker 工作池  这个动作只能发生一次，一个zinx 框架只能有一个worker池
func (m *MsgHandler) StartWorkerPool() {
	// 根据WorkerPoolSize 分别开启Worker， 每个Worker用一个goroutine来承载
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		// 一个worker 被启动
		// 给当前worker 对应的channel 消息队列开辟空间, 第 i 个 worker 对应第 i 个 channel
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerQueueLength)
		// 启动当前的worker, 阻塞等待消息从channel 传过来
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

// 启动一个worker 工作流程
func (m *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Printf("Worker ID = %d is started ...\n", workerID)

	// 阻塞等待对应的消息队列
	for {
		select {
		// 如果有消息过来， 出列的就是一个客户端的Request, 执行当前Request所绑定的业务
		case request := <-taskQueue:
			m.DoMsgHandle(request)

		}
	}
}

// 将消息交给TaskQueue， 由Worker进行处理
func (m *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1 将消息平均分配给不同worker
	// 根据客户端建立的ConnID 来进行分配
	workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		" request MsgID = ", request.GetMsgID(), " to worker ID = ", workerID)
	// 2 将消息发送给对应的worker的TaskQueue即可
	m.TaskQueue[workerID] <- request
}
