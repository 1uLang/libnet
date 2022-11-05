package message

import (
	"encoding/binary"
	"errors"
	"github.com/1uLang/libnet/message"
	"sync/atomic"
)

const (
	MessageVersion = 0x1 // 版本， 0-f

	MessageIdIndex      = 1
	MessageTypeIndex    = 9
	MessageLengthIndex  = 10
	MessageHeaderLength = 14
)

// Message tls 通讯消息体
// Type：
// 0x00 AH 登录信息包  AH ----> 控制器
// 0x01 控制器 AH登录响应包  控制器 ------> AH
// 0x02 AH 登出响应包 AH  --------> 控制器
// 0x03 AH/控制器 keepalive 心跳包
// 0x04 AH AH服务信息包

const (
	LoginRequestCode         = 0x00 //ah/ih ----> control 登录消息
	LoginResponseCode        = 0x01 //control ----> ah/ih 登录响应消息
	AHLogoutRequestCode      = 0x02 // ah ----> control 注销消息
	KeepaliveRequestCode     = 0x03 // ah <----> control 心跳消息
	ServerProtectRequestCode = 0x04 // control ----> ah ah保护服务消息
	IHOnlineRequestCode      = 0x05 // control ----> ah ih认证消息
	AHListRequestCode        = 0x06 // control ----> ih ah信息列表
	IHLogoutRequestCode      = 0x07 // ih ----> control 注销请求消息
	IHOnlineResponseCode     = 0x08 // ah ----> control ih上线后ah业务相关数据信息体响应消息
	CustomRequestCode        = 0xff // 自定义消息
)

var (
	ErrBufferInvalidIsNil  = errors.New("buffer: invalid is nil ")
	ErrBufferInvalidStart  = errors.New("buffer: invalid start byte")
	ErrBufferInvalidHeader = errors.New("buffer: invalid header")
	ErrBufferDataTooLong   = errors.New("buffer: data too long bytes")
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

// Marshal 编码消息
func (this *Message) Marshal() []byte {

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
	return result
}

// MsgId 获取消息ID
func (this *Message) MsgId() uint64 {
	return this.Id
}

// SetData 设置消息体
func (this *Message) SetData(buf []byte) {
	if this.Data == nil {
		this.Data = buf
	} else {
		this.Data = append(this.Data, buf...)
	}
}

// HeaderLength 获取消息头长度
func (this *Message) HeaderLength() uint32 {
	return MessageHeaderLength
}

// GetLength 获取消息体长度
func (this *Message) GetLength() uint32 {
	return this.Length
}

// CheckHeader 分析并检测消息头
func CheckHeader(buf []byte) (message.MessageI, error) {
	msg := Message{}
	if buf == nil && len(buf) == 0 {
		return nil, ErrBufferInvalidIsNil
	}
	msg.Version = buf[0]
	// 检查消息版本
	if len(buf) > 0 && buf[0] != MessageVersion {
		return nil, ErrBufferInvalidStart
	}

	if len(buf) < MessageHeaderLength {
		return nil, ErrBufferInvalidHeader
	}

	l := binary.BigEndian.Uint32(buf[MessageLengthIndex : MessageLengthIndex+4])
	if l > message.MaxBufferSize { // 每次通讯数据不超过一定尺寸
		return nil, ErrBufferDataTooLong
	}
	msg.Length = l
	// 解析Header
	msg.Id = binary.BigEndian.Uint64(buf[MessageIdIndex : MessageIdIndex+8])
	msg.Type = buf[MessageTypeIndex]
	return &msg, nil
}
