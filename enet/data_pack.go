package enet

import (
	"EagleNet/configs"
	"EagleNet/eiface"
	"bytes"
	"encoding/binary"
	"fmt"
)

/* TLV格式消息封装
msgLen msgID xxxxxxxxxxxxxxx
<- header -> <-   body    ->
*/

//消息TLV格式封包/解包
type DataPack struct {
	sizeMsgLen    uint32 //消息长度所占字节
	sizeMsgID     uint32 //消息id所占字节
	sizeMsgHeader uint32 //消息头所占字节
	maxPkgSize    int    //消息最大长度
}

const (
	defaultSizeMsgLen = 4
	defaultSizeMsgID  = 4
)

func NewDataPack(dataPackCfg *configs.DataPack) *DataPack {
	return &DataPack{
		sizeMsgLen:    defaultSizeMsgLen,
		sizeMsgID:     defaultSizeMsgID,
		sizeMsgHeader: defaultSizeMsgLen + defaultSizeMsgID,
		maxPkgSize:    dataPackCfg.MaxPkgSize,
	}
}

func (d *DataPack) GetHeaderLen() uint32 {
	return d.sizeMsgHeader
}

// TLV格式封包
func (d *DataPack) Pack(msg eiface.IMessage) ([]byte, error) {
	var err error
	buff := bytes.NewBuffer([]byte{})

	if err = binary.Write(buff, binary.LittleEndian, msg.GetLen()); err != nil {
		return nil, err
	}
	if err = binary.Write(buff, binary.LittleEndian, msg.GetID()); err != nil {
		return nil, err
	}
	if err = binary.Write(buff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// TLV格式解包,只解header部分
func (d *DataPack) UnPack(header []byte) (eiface.IMessage, error) {
	var err error
	buff := bytes.NewReader(header)
	msg := &Message{}

	if err = binary.Read(buff, binary.LittleEndian, &msg.Len); err != nil {
		return nil, err
	}
	if d.maxPkgSize > 0 && int(msg.Len) > d.maxPkgSize {
		return nil, fmt.Errorf("message length(%d) exceed MaxPackageSize(%d)", msg.Len, d.maxPkgSize)
	}

	if err = binary.Read(buff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}

	// 只解header部分,body部分数据需要在外面再读一次,读取长度为msg.Len的数据

	return msg, nil
}
