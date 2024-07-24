package eiface

//连接管理接口
type IConnManager interface {
	//添加连接
	Add(conn IConnection)
	//删除连接
	Remove(connID uint32)
	//获取连接
	Get(connID uint32) (IConnection, error)
	//获取连接总数
	Len() int
	//清除所有连接
	Clear()
}
