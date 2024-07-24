package eiface

//消息处理接口
type IMsgHandler interface {
	//为消息注册路由
	AddRouter(msgID uint32, router IRouter)
	//查找并执行消息的路由(同步)
	DoMsgHandler(req IRequest)
	//打印所有路由
	PrintAllRouters()
	//启动Worker工作池
	StartWorkerPool()
	//关闭Worker工作池
	StopWorkerPool()
	//异步处理请求
	ProcessRequestAsync(req IRequest)
}
