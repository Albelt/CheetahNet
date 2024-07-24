package enet

import "EagleNet/eiface"

//IRequest的实现
type Request struct {
	//请求关联的连接
	conn eiface.IConnection
	//请求数据
	msg eiface.IMessage
}

func (r *Request) GetConnection() eiface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetID()
}
