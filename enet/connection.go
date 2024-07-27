package enet

import (
	"EagleNet/configs"
	"EagleNet/eiface"
	"EagleNet/pkg/log"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

// IConnection实现
type Connection struct {
	//tcp连接
	Conn *net.TCPConn
	//连接id
	ConnID uint32
	//是否已关闭
	Closed bool
	//告知当前连接已退出
	ExitChan chan bool
	//Reader和Writer之间通信的channel(无缓冲,同步读写)
	msgChann chan []byte
	//关联的服务器
	srv eiface.IServer

	//用户自定义的连接属性
	properties map[string]interface{}
	pMutex     sync.RWMutex

	//配置
	dataPackCfg *configs.DataPack
}

func NewConnection(srv eiface.IServer, conn *net.TCPConn, connID uint32, dataPackCfg *configs.DataPack) eiface.IConnection {
	c := &Connection{
		srv:         srv,
		Conn:        conn,
		ConnID:      connID,
		Closed:      false,
		ExitChan:    make(chan bool),
		msgChann:    make(chan []byte),
		properties:  make(map[string]interface{}),
		dataPackCfg: dataPackCfg,
	}

	//将连接添加到连接管理器
	c.srv.GetConnManager().Add(c)
	return c
}

func (c *Connection) Start() {
	log.Infof("Connection Start, connID:%d, remoteAddr:%s", c.ConnID, c.GetRemoteAddr().String())

	//启动读协程
	go c.startReader()

	//启动写协程
	go c.startWriter()

	//调用钩子函数
	c.srv.CallHookOnConnStart(c)
}

//读数据的协程
func (c *Connection) startReader() {
	log.Infof("ConnID:%d, Reader goroutine is running...", c.GetConnID())
	defer func() {
		log.Infof("ConnID:%d, Reader is exit", c.GetConnID())
		c.Stop()
	}()

	// 业务逻辑
	var err error
	dp := NewDataPack(c.dataPackCfg)
	for {
		//读取固定长度的header,解析为Message
		header := make([]byte, dp.GetHeaderLen())
		_, err = io.ReadFull(c.GetTCPConn(), header)
		if err != nil {
			break
		}
		msg, err := dp.UnPack(header)
		if err != nil {
			break
		}

		//读取指定长度的body
		if msg.GetLen() == 0 {
			continue
		}
		body := make([]byte, msg.GetLen())
		if _, err := io.ReadFull(c.GetTCPConn(), body); err != nil {
			break
		}
		msg.SetData(body)

		//将连接和数据放到Request中,然后调用消息处理器的处理方法
		req := &Request{
			conn: c,
			msg:  msg,
		}
		c.srv.GetMsgHandler().ProcessRequestAsync(req)
	}

	log.Infof("Reader got err:%s", err.Error())
}

//写数据的协程
func (c *Connection) startWriter() {
	log.Infof("ConnID:%d, Writer goroutine is running...", c.GetConnID())
	defer func() {
		log.Infof("ConnID:%d, Writer is exit", c.GetConnID())
	}()

	//阻塞等待channel的消息,往客户端发送数据
	var err error
Loop:
	for {
		select {
		case data := <-c.msgChann:
			if _, err = c.Conn.Write(data); err != nil {
				log.Infof("Writer got err:%s", err.Error())
			}
		case <-c.ExitChan:
			break Loop
		}
	}
}

func (c *Connection) Stop() {
	log.Infof("Connection Stop, connID:%d, remoteAddr:%s", c.ConnID, c.GetRemoteAddr().String())

	// 判断关闭标志
	if c.Closed {
		return
	}
	c.Closed = true

	//调用钩子函数
	c.srv.CallHookOnConnStop(c)

	//关闭连接
	c.Conn.Close()

	//告知Writer退出
	c.ExitChan <- true

	//关闭管道
	//关闭无缓冲的管道后,读管道会立即返回默认值,写管道会panic
	close(c.ExitChan)
	close(c.msgChann)

	//连接管理器删除该连接
	c.srv.GetConnManager().Remove(c.ConnID)
}

func (c *Connection) GetTCPConn() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(data []byte, msgID uint32) error {
	if c.Closed {
		return errors.New("Connection is closed.")
	}

	//将消息进行封装
	msg := &Message{
		ID:   msgID,
		Len:  uint32(len(data)),
		Data: data,
	}

	dp := NewDataPack(c.dataPackCfg)
	bytes, err := dp.Pack(msg)
	if err != nil {
		return fmt.Errorf("SendMsg.Pack msg err:%s", err.Error())
	}

	//发送消息
	//_, err = c.Conn.Write(bytes)
	//if err != nil {
	//	return fmt.Errorf("SendMsg.Write msg err:%s", err.Error())
	//}

	//待发送数据写入管道,由Writer发送给客户端
	c.msgChann <- bytes

	return nil
}

func (c *Connection) SetProperty(k string, v interface{}) {
	c.pMutex.Lock()
	defer c.pMutex.Unlock()

	c.properties[k] = v
}

func (c *Connection) GetProperty(k string) (interface{}, error) {
	c.pMutex.RLock()
	defer c.pMutex.RUnlock()

	if v, ok := c.properties[k]; ok {
		return v, nil
	} else {
		return nil, errors.New("property not exist")
	}
}

func (c *Connection) DelProperty(k string) {
	c.pMutex.Lock()
	defer c.pMutex.Unlock()

	delete(c.properties, k)
}
