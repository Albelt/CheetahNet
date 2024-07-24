package eiface

// 服务器接口
type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()
	//注册一个路由
	AddRouter(msgID uint32, router IRouter)
	//获取消息管理器
	GetMsgHandler() IMsgHandler
	//获取连接管理器
	GetConnManager() IConnManager
	//设置钩子函数(连接创建之后)
	SetHookOnConnStart(hook ConnHook)
	//设置钩子函数(连接关闭之前)
	SetHookOnConnStop(hook ConnHook)
	//调用钩子函数(连接创建之后)
	CallHookOnConnStart(conn IConnection)
	//调用钩子函数(连接关闭之前)
	CallHookOnConnStop(conn IConnection)
}

//连接事件的钩子函数(在连接创建/关闭时执行特定的逻辑)
type ConnHook func(conn IConnection)
