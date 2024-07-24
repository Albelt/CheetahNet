package enet

import (
	"EagleNet/configs"
	"EagleNet/eiface"
	"fmt"
)

// IMsgHandler实现
type MsgHandler struct {
	//MsgID到Router的映射
	routerMap map[uint32]eiface.IRouter
	//协程池Worker数量
	workerPoolSize int
	//Worker绑定的队列(每个Worker一个)
	workerQueues []chan eiface.IRequest
	//关闭Worker池的信号
	workerCloseSig chan struct{}
	//配置
	workerPoolCfg *configs.WorkerPool
}

func NewMsgHandler(workerPoolCfg *configs.WorkerPool) *MsgHandler {
	return &MsgHandler{
		routerMap:      make(map[uint32]eiface.IRouter),
		workerPoolCfg:  workerPoolCfg,
		workerPoolSize: workerPoolCfg.PoolSize,
		workerQueues:   make([]chan eiface.IRequest, workerPoolCfg.PoolSize),
	}
}

func (h *MsgHandler) AddRouter(msgID uint32, router eiface.IRouter) {
	if _, ok := h.routerMap[msgID]; ok {
		panic(fmt.Sprintf("Duplicated router add to msgID(%d)", msgID))
	}

	h.routerMap[msgID] = router
}

//查找并执行消息的路由
func (h *MsgHandler) DoMsgHandler(req eiface.IRequest) {
	//查询msgID对应的路由
	router, ok := h.routerMap[req.GetMsgID()]
	if !ok {
		fmt.Printf("MsgID(%d)'s router not found\n", req.GetMsgID())
		return
	}

	//调用路由方法
	router.PreHandler(req)
	router.Handler(req)
	router.PostHandler(req)
}

func (h *MsgHandler) PrintAllRouters() {
	if len(h.routerMap) == 0 {
		fmt.Println("No routers")
		return
	}

	fmt.Println("All Routers:")
	for msgID, router := range h.routerMap {
		fmt.Printf("MsgID(%d) -> Router(%T)\n", msgID, router)
	}
	fmt.Println()
}

func (h *MsgHandler) StartWorkerPool() {
	for i := 0; i < h.workerPoolSize; i++ {
		//初始化worker队列
		h.workerQueues[i] = make(chan eiface.IRequest, h.workerPoolCfg.QueueSize)

		//启动worker
		go h.startWorker(i, h.workerQueues[i])
	}
}

//启动单个Worker
func (h *MsgHandler) startWorker(workerID int, taskQ chan eiface.IRequest) {
	fmt.Printf("Worker(ID=%d) starting...\n", workerID)
	defer func() {
		fmt.Printf("Worker(ID=%d) exit\n", workerID)
	}()

	for {
		select {
		case req := <-taskQ: //收到数据,调用路由进行处理
			h.DoMsgHandler(req)
		case <-h.workerCloseSig: //收到退出信号,Worker退出
			return
		}
	}
}

func (h *MsgHandler) StopWorkerPool() {
	close(h.workerCloseSig)
}

//处理请求,将请求写入Worker的队列即可
func (h *MsgHandler) ProcessRequestAsync(req eiface.IRequest) {
	//将请求平均分配给各个Worker: round-robin算法, 根据连接ID做哈希
	workerID := int(req.GetConnection().GetConnID() % uint32(h.workerPoolSize))

	//若worker的队列满了,写数据会阻塞
	h.workerQueues[workerID] <- req
}
