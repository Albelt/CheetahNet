package configs

type Config struct {
	Server     *Server     `json:"server"`
	WorkerPool *WorkerPool `json:"worker_pool"`
	DataPack   *DataPack   `json:"data_pack"`
}

type Server struct {
	Name    string
	IP      string
	Port    int
	MaxConn int
}

type WorkerPool struct {
	PoolSize  int
	QueueSize int
}

type DataPack struct {
	MaxPkgSize int
}
