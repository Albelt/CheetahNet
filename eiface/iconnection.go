package eiface

import "net"

// 客户端连接接口
type IConnection interface {
	//启动连接,开始处理业务
	Start()
	//结束连接
	Stop()
	//获取底层tcp连接
	GetTCPConn() *net.TCPConn
	//获取连接id
	GetConnID() uint32
	//获取客户端的地址
	GetRemoteAddr() net.Addr
	//发送数据给客户端
	SendMsg(data []byte, msgID uint32) error
	//设置连接属性
	SetProperty(k string, v interface{})
	//获取连接属性
	GetProperty(k string) (interface{}, error)
	//删除连接属性
	DelProperty(k string)
}

// 处理连接业务的方法
// @param: tcp连接, 读取数据, 读取数据长度
// @return: error
type HandleFunc func(*net.TCPConn, []byte, int) error
