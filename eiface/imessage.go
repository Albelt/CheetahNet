package eiface

//将请求的消息封装
type IMessage interface {
	GetID() uint32   //获取消息id
	GetLen() uint32  //获取消息长度
	GetData() []byte //获取消息内容

	SetID(uint32)   //设置消息id
	SetLen(uint32)  //设置消息长度
	SetData([]byte) //设置消息内容
}
