package message

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/1uLang/libnet/utils/maps"
	"sync/atomic"
)

const (
	MessageVersion = 0x1 // 版本， 0-f

	MessageIdIndex      = 1
	MessageTypeIndex    = 9
	MessageLengthIndex  = 10
	MessageHeaderLength = 14
)

var (
	messageId = uint64(0)
)

type Message struct {
	Version byte
	Id      uint64
	Type    byte
	Length  uint32
	Data    []byte
}

// 解析选项
// 选项格式：[4字节选项内容长度] [ 选项内容 ] [数据]
func (this *Message) DecodeOptions() (maps.Map, error) {
	if len(this.Data) < 4 {
		return nil, nil
	}

	length := binary.BigEndian.Uint32(this.Data[:4])
	if length == 0 {
		return nil, nil
	}

	data := this.Data[4 : 4+length]
	options := maps.Map{}
	err := json.Unmarshal(data, &options)

	this.Data = this.Data[4+length:]

	return options, err
}

// 编码消息
func (this *Message) Marshal() []byte {
	fmt.Println(this.Data)
	if this.Id <= 0 {
		this.Id = atomic.AddUint64(&messageId, 1)
	}

	if this.Version == 0 {
		this.Version = MessageVersion
	}
	result := []byte{this.Version}

	// ID
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, this.Id)
	result = append(result, buf...)

	// Type
	result = append(result, this.Type)

	// Length
	this.Length = uint32(len(this.Data))
	binary.BigEndian.PutUint32(buf, this.Length)
	result = append(result, buf[:4]...)

	// Data
	result = append(result, this.Data...)
	fmt.Println("====", this.Length, this.Data)
	return result
}
