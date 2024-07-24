package configs

import (
	"encoding/json"
	"io/ioutil"
)

//从json文件中加载配置,没有配置文件则使用默认配置
func LoadConfigs(path string) (*Config, error) {
	//默认配置
	cfg := Config{
		Server: &Server{
			Name:    "EagleNet",
			IP:      "0.0.0.0",
			Port:    8888,
			MaxConn: 1000,
		},
		WorkerPool: &WorkerPool{
			PoolSize:  10,
			QueueSize: 100,
		},
		DataPack: &DataPack{MaxPkgSize: 2048},
	}

	if path == "" {
		return &cfg, nil
	}

	//从文件中加载配置
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buff, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
