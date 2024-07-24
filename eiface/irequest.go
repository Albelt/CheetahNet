package eiface

//客户端请求的连接和数据封装
type IRequest interface {
	//获取当前连接
	GetConnection() IConnection
	//获取请求数据
	GetData() []byte
	//获取消息id
	GetMsgID() uint32
}
