package message

import (
	"encoding/binary"
	"errors"
	"strconv"
)

const MaxBufferSize = 2 * 1024 * 1024 // 最大消息长度

var (
	ErrBufferInvalidStart = errors.New("buffer: invalid start byte")
	ErrBufferInvalidId    = errors.New("buffer: invalid id")
)

type Buffer struct {
	OptValidateId bool // 防重放攻击开关

	buf    []byte
	msgId  uint64
	msgLen int
	msg    *Message

	onMessage func(msg *Message)
	onError   func(err error) //todo:通过该函数的触发 记录错误ip并记录，次数超过阈值 在IPset中拉黑 并将该ip上传到admin

	hasError bool
}

func NewBuffer() *Buffer {
	return &Buffer{}
}

//	[version]     [ id ]   [type]    [length]    [data]
//
// [1字节版本号] [8位请求ID][1字节数据类型][4字节数据长度][数据]
func (this *Buffer) Write(rawBuf []byte) {
	if this.hasError {
		return
	}

	if len(rawBuf) == 0 {
		return
	}

	buf := make([]byte, len(rawBuf))
	copy(buf, rawBuf)

	if this.msgLen > 0 {
		buf = this.writeBytes(buf)
		if len(buf) == 0 {
			return
		}
	}

	if len(this.buf) == 0 {
		this.buf = make([]byte, len(buf))
		copy(this.buf, buf)
	} else {
		this.buf = append(this.buf, buf...)
	}

	for {
		// 检查消息版本
		if len(this.buf) > 0 && this.buf[0] != MessageVersion {
			if this.onError != nil {
				this.reset()
				this.onError(ErrBufferInvalidStart)
			}
			return
		}

		if len(this.buf) < MessageHeaderLength {
			return
		}

		l := binary.BigEndian.Uint32(this.buf[MessageLengthIndex : MessageLengthIndex+4])
		if l > MaxBufferSize { // 每次通讯数据不超过一定尺寸
			if this.onError != nil {
				this.onError(errors.New("data too long '" + strconv.Itoa(int(l)) + "' bytes"))
			}
			return
		}

		// 解析Header
		msgId := binary.BigEndian.Uint64(this.buf[MessageIdIndex : MessageIdIndex+8])
		dataType := this.buf[MessageTypeIndex]

		// 防重放攻击
		if this.OptValidateId && msgId <= this.msgId {
			if this.onError != nil {
				this.reset()
				this.onError(ErrBufferInvalidId)
			}
			return
		}
		this.msgId = msgId
		this.msg = nil

		msg := new(Message)
		msg.Id = msgId
		msg.Version = this.buf[0]
		msg.Type = dataType
		msg.Length = l

		if l == 0 {
			if this.onMessage != nil {
				this.onMessage(msg)

				// 由于onMessage可能会改变buffer，所以这里需要做判断
				if len(this.buf) < MessageHeaderLength {
					return
				}
			}
		} else {
			this.msg = msg
		}

		this.msgLen = int(l)

		// 写入剩下的数据
		this.buf = this.buf[MessageHeaderLength:]
		if this.msgLen > 0 {
			this.buf = this.writeBytes(this.buf)
		} else {
			if this.msg != nil && this.onMessage != nil {
				this.onMessage(this.msg)
			}
		}
	}
}

func (this *Buffer) OnMessage(f func(msg *Message)) {
	this.onMessage = f
}

func (this *Buffer) OnError(f func(err error)) {
	this.onError = f
}

func (this *Buffer) Reset() {
	this.reset()
}
func (this *Buffer) reset() {
	// 可能有异步操作的风险，暂时先注释掉
	return

	this.buf = nil
	this.msgLen = 0
	this.msg = nil
	this.msgId = 0
	this.hasError = false
}

func (this *Buffer) writeBytes(buf []byte) []byte {
	l := len(buf)
	if l <= this.msgLen {
		this.msgLen = this.msgLen - l

		if this.msg != nil && this.onMessage != nil {
			if this.msg.Data == nil {
				this.msg.Data = buf
			} else {
				this.msg.Data = append(this.msg.Data, buf...)
			}
			if this.msgLen == 0 {
				this.onMessage(this.msg)
			}
		}

		return nil
	}

	// if l > msgLen
	if this.msg != nil && this.onMessage != nil {
		if this.msg.Data == nil {
			this.msg.Data = buf[:this.msgLen]
		} else {
			this.msg.Data = append(this.msg.Data, buf[:this.msgLen]...)
		}
		this.onMessage(this.msg)
	}

	buf = buf[this.msgLen:]
	this.msgLen = 0

	return buf
}
