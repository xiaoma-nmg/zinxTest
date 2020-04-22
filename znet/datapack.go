package znet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"zinx/utils"
	"zinx/ziface"
)

// 解决tcp 粘包问题的封包拆包的模块

// TLV: type length value 格式的消息格式

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包长度的方法
func (d *DataPack) GetHeadLen() uint32 {
	return 8 // DataLen 4KB  ID 4KB
}

// 封包方法
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 存放byte字节流的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将dataLen 写进databuff
	err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLength())
	if err != nil {
		return nil, err
	}
	// 将msgId 写进databuff
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}
	// 将data数据 写进databuff
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法
func (d *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入数据读取的 ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压head信息， 得到dataLen 和 MsgID
	msg := &Message{}

	// 读dataLen
	err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}

	// 读MsgID
	err = binary.Read(dataBuff, binary.LittleEndian, &msg.Id)
	if err != nil {
		return nil, err
	}

	//判断是否已经超出了 允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large package")
	}
	return msg, nil
}
