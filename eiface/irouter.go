package eiface

//路由接口
type IRouter interface {
	//处理业务之前的方法
	PreHandler(req IRequest)
	//处理业务的方法
	Handler(req IRequest)
	//处理业务之后的方法
	PostHandler(req IRequest)
}
