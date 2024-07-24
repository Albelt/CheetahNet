package eiface

//对消息进行封包/拆包
type IDataPack interface {
	GetHeaderLen() uint32              //获取header长度
	Pack(msg IMessage) ([]byte, error) //封包
	UnPack([]byte) (IMessage, error)   //拆包
}
