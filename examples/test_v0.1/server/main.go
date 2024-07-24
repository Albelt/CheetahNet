package main

import (
	"EagleNet/configs"
	"EagleNet/eiface"
	"EagleNet/enet"
	"fmt"
)

func main() {
	//加载配置
	cfg, err := configs.LoadConfigs("")
	if err != nil {
		panic(err)
	}

	//创建服务器
	srv := enet.NewServer(cfg)

	//添加router
	srv.AddRouter(msgID_01, &myRouter{})

	//设置连接钩子函数
	srv.SetHookOnConnStart(func(conn eiface.IConnection) {
		_ = conn.SendMsg([]byte("Connection begin."), msgID_01)
	})

	//启动服务
	srv.Serve()
}

const (
	msgID_01 = uint32(1)
)

type myRouter struct {
	enet.BaseRouter
}

func newMyRouter() *myRouter {
	return &myRouter{BaseRouter: enet.BaseRouter{}}
}

func (r *myRouter) Handler(req eiface.IRequest) {
	fmt.Printf("myRouter.Handler received msg, ID:%d, data:%s\n", req.GetMsgID(), req.GetData())

	sayHello := "Hello, this is EagleNet!"
	err := req.GetConnection().SendMsg([]byte(sayHello), msgID_01)
	if err != nil {
		fmt.Printf("Handler err:%s\n", err.Error())
	}
}