package enet

import (
	"EagleNet/configs"
	"EagleNet/eiface"
	"fmt"
	"net"
	"EagleNet/pkg/log"
)

// IServer的实现
type Server struct {
	Name       string //服务器名称
	IPVersion  string //ip版本
	IP         string
	Port       int
	msgHandler eiface.IMsgHandler  //消息管理模块
	connMgr    eiface.IConnManager //连接管理模块

	//连接的钩子函数
	connHookOnConnStart eiface.ConnHook
	connHookOnConnStop  eiface.ConnHook

	//配置
	serverCfg   *configs.Server
	dataPackCfg *configs.DataPack
}

func NewServer(cfg *configs.Config) eiface.IServer {
	return &Server{
		Name:        cfg.Server.Name,
		IPVersion:   "tcp4",
		IP:          cfg.Server.IP,
		Port:        cfg.Server.Port,
		msgHandler:  NewMsgHandler(cfg.WorkerPool),
		connMgr:     NewConnManager(),
		serverCfg:   cfg.Server,
		dataPackCfg: cfg.DataPack,
	}
}

func (s *Server) Start() {
	go s.start()
}

func (s *Server) start() {
	// 解析地址
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		log.Errorf("ResolveTCPAddr err:%s", err.Error())
		return
	}

	// 创建socket并监听
	lis, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		log.Errorf("ListenTCP err:%s", err.Error())
		return
	}
	defer lis.Close()

	log.Infof("[Start] Server(%s) on %s: %s:%d", s.Name, s.IPVersion, s.IP, s.Port)

	// 接收客户端连接,处理读写
	var connID uint32
	for {
		conn, err := lis.AcceptTCP()
		if err != nil {
			log.Errorf("AcceptTCP err:%s", err.Error())
			break
		}
		connID++

		//检测连接数据是否超过限制
		if s.connMgr.Len() >= s.serverCfg.MaxConn {
			//TODO:给客户端响应一个错误
			conn.Close()
			log.Warnf("Connection exceed max number %d", s.serverCfg.MaxConn)
			continue
		}

		//创建zinx连接并添加到连接管理器
		zinxConn := NewConnection(s, conn, connID, s.dataPackCfg)
		//启动连接的业务处理
		go zinxConn.Start()
	}
}

func (s *Server) Stop() {
	//回收资源
	s.msgHandler.StopWorkerPool()
	s.connMgr.Clear()
}

func (s *Server) Serve() {
	//打印所有路由
	s.msgHandler.PrintAllRouters()

	//启动工作池
	s.msgHandler.StartWorkerPool()

	//异步启动服务,不会阻塞
	s.Start()

	//在这里阻塞
	//TODO:检测退出信号,调用Stop()
	select {}
}

func (s *Server) GetMsgHandler() eiface.IMsgHandler {
	return s.msgHandler
}

func (s *Server) GetConnManager() eiface.IConnManager {
	return s.connMgr
}

func (s *Server) AddRouter(msgID uint32, router eiface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

func (s *Server) SetHookOnConnStart(hook eiface.ConnHook) {
	s.connHookOnConnStart = hook
}

func (s *Server) SetHookOnConnStop(hook eiface.ConnHook) {
	s.connHookOnConnStop = hook
}

func (s *Server) CallHookOnConnStart(conn eiface.IConnection) {
	if s.connHookOnConnStart != nil {
		s.connHookOnConnStart(conn)
	}
}

func (s *Server) CallHookOnConnStop(conn eiface.IConnection) {
	if s.connHookOnConnStop != nil {
		s.connHookOnConnStop(conn)
	}
}
