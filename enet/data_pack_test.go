package enet

import (
	"EagleNet/configs"
	"io"
	"net"
	"reflect"
	"testing"
	"time"
)

var (
	dpConfig = &configs.DataPack{MaxPkgSize: 2048}
)

//测试封包解包
func TestDataPack(t *testing.T) {
	dp := NewDataPack(dpConfig)

	msg := &Message{
		ID:   1,
		Len:  5,
		Data: []byte{'h', 'e', 'l', 'l', 'o'},
	}

	packedMsg, err := dp.Pack(msg)
	if err != nil {
		t.Fatal(err)
	}

	unPackedMsg, err := dp.UnPack(packedMsg)
	if err != nil {
		t.Fatal(err)
	}

	if unPackedMsg.GetID() != msg.GetID() || unPackedMsg.GetLen() != msg.GetLen() ||
		!reflect.DeepEqual(msg.GetData(), unPackedMsg.GetData()) {
		t.Fatalf("Expected:%+vGot:%+v", msg, unPackedMsg)
	}
}

//测试TCP粘包拆包
func TestTCPPacketStickingAndSpliting(t *testing.T) {
	const (
		srvAddr = "127.0.0.1:9999"
	)

	//启动服务器
	lis, err := net.Listen("tcp", srvAddr)
	if err != nil {
		return
	}
	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				return
			}

			go func() {
				for {
					//读取并解包header
					dp := NewDataPack(dpConfig)
					header := make([]byte, dp.GetHeaderLen())
					_, err := io.ReadFull(conn, header)
					if err != nil {
						return
					}
					unPackData, err := dp.UnPack(header)
					if err != nil {
						return
					}

					//根据header中的length,读body
					msg := unPackData.(*Message)
					msg.Data = make([]byte, msg.GetLen())
					_, err = io.ReadFull(conn, msg.Data)
					if err != nil {
						return
					}

					t.Logf("Recv msg, id:%d, len:%d, msg:%v", msg.GetID(), msg.GetLen(), msg.GetData())
				}
			}()
		}
	}()

	//启动客户端
	go func() {
		conn, err := net.Dial("tcp", srvAddr)
		if err != nil {
			return
		}

		dp := NewDataPack(dpConfig)

		//封装2个数据包,组成粘包
		msg1 := &Message{
			ID:   1,
			Len:  5,
			Data: []byte("hello"),
		}
		msg2 := &Message{
			ID:   2,
			Len:  3,
			Data: []byte("bye"),
		}

		sendData, err := dp.Pack(msg1)
		if err != nil {
			return
		}
		tmp, err := dp.Pack(msg2)
		if err != nil {
			return
		}
		sendData = append(sendData, tmp...)

		_, err = conn.Write(sendData)
		if err != nil {
			return
		}

		t.Logf("Cli send msg:%v", sendData)
		time.Sleep(time.Second)
		
		//将1个数据包分两次发送,测试拆包
		msg3 := &Message{
			ID:   0,
			Len:  10,
			Data: []byte("helloworld"),
		}
		sendData3, err := dp.Pack(msg3)
		_, _ = conn.Write(sendData3[:5])
		_, _ = conn.Write(sendData3[5:])
		t.Logf("Cli send msg:%v", sendData3)
	}()

	//阻塞
	select {
	case <-time.After(time.Second * 2):
		t.Logf("finished.")
	}
}
