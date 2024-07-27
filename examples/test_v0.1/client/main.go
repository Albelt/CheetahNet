package main

import (
	"EagleNet/configs"
	"EagleNet/pkg/log"
	"errors"
	"fmt"
	"io"
	"EagleNet/enet"
	"net"
)

// 模拟客户端

const (
	msgID_01 = uint32(1)
)

func main() {
	//加载配置
	cfg, err := configs.LoadConfigs("")
	if err != nil {
		panic(err)
	}

	// 连接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", cfg.Server.Port))
	if err != nil {
		panic(err)
	}

	// 读写数据,使用znet提供的封包/解包方法
	dp := enet.NewDataPack(cfg.DataPack)

	// 读数据
	go func() {
		for {
			header := make([]byte, dp.GetHeaderLen())
			_, err = io.ReadFull(conn, header)
			if err != nil {
				log.Infof("Read err:%s", err.Error())
				if errors.Is(err, io.EOF) {
					break
				} else {
					continue
				}
			}
			unPackdMsg, err := dp.UnPack(header)
			if err != nil {
				log.Infof("Unpack err:%s", err.Error())
				continue
			}
			if unPackdMsg.GetLen() > 0 {
				body := make([]byte, unPackdMsg.GetLen())
				_, _ = io.ReadFull(conn, body)
				unPackdMsg.SetData(body)
			}
			log.Infof("read msg, ID:%d, data:%s", unPackdMsg.GetID(), unPackdMsg.GetData())
		}
	}()

	//写数据
	go func() {
		var input string
		for {
			_, err := fmt.Scanf("%s", &input)
			if err != nil {
				log.Infof("Read err:%s", err)
				continue
			}

			msg := &enet.Message{
				ID:   msgID_01,
				Len:  uint32(len(input)),
				Data: []byte(input),
			}
			bytes, _ := dp.Pack(msg)

			// 写数据
			_, err = conn.Write(bytes)
			if err != nil {
				log.Infof("Write err:%s", err.Error())
				break
			}
			log.Infof("write msg, ID:%d, data:%s", msg.GetID(), msg.GetData())
		}
	}()

	//主协程阻塞
	select {}
}
