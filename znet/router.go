package znet

import "zinx/ziface"

// 实现router 时， 先嵌入 BaseRouter 基类， 然后根据需要对这些方法重写
type BaseRouter struct{}

//BaseRouter 的方法都为空，目的是为了能够覆盖 接口。它的子结构可以只实现其中任意一个方法
//在处理 conn 业务之前的钩子方法 hook
func (b *BaseRouter) PreHandle(ziface.IRequest) {}

//在处理 conn 业务的方法
func (b *BaseRouter) Handle(ziface.IRequest) {}

//在处理 conn 业务之后的方法
func (b *BaseRouter) PostHandle(ziface.IRequest) {}
