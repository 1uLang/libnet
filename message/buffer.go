package message

import (
	"errors"
)

const MaxBufferSize = 2 * 1024 * 1024 // 最大消息长度

var (
	ErrBufferInvalidId = errors.New("buffer: invalid id")
)

type Buffer struct {
	OptValidateId bool // 防重放攻击开关

	buf    []byte
	msgId  uint64
	msgLen uint32
	msg    MessageI

	onMessage  func(msg MessageI)
	onError    func(err error)
	parserFunc func([]byte) (MessageI, error)
	hasError   bool
}

func NewBuffer(parserFunc func([]byte) (MessageI, error)) *Buffer {
	return &Buffer{
		parserFunc: parserFunc,
	}
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
		// 检查消息头
		msg, err := this.parserFunc(this.buf)
		if err != nil {
			if this.onError != nil {
				this.reset()
				this.onError(err)
			}
			return
		}
		// 防重放攻击
		if this.OptValidateId && msg.MsgId() <= this.msgId {
			if this.onError != nil {
				this.reset()
				this.onError(ErrBufferInvalidId)
			}
			return
		}
		this.msgId = msg.MsgId()
		if msg.GetLength() == 0 {
			if this.onMessage != nil {
				this.onMessage(msg)
				// 由于onMessage可能会改变buffer，所以这里需要做判断
				if len(this.buf) < int(msg.HeaderLength()) {
					return
				}
			}
		} else {
			this.msg = msg
		}

		this.msgLen = msg.GetLength()

		// 写入剩下的数据
		this.buf = this.buf[msg.HeaderLength():]
		if this.msgLen > 0 {
			this.buf = this.writeBytes(this.buf)
		} else {
			if this.msg != nil && this.onMessage != nil {
				this.onMessage(this.msg)
			}
		}
	}
}

func (this *Buffer) OnMessage(f func(msg MessageI)) {
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
}

func (this *Buffer) writeBytes(buf []byte) []byte {
	l := uint32(len(buf))
	if l <= this.msgLen {
		this.msgLen = this.msgLen - l

		if this.msg != nil && this.onMessage != nil {
			this.msg.SetData(buf)
			if this.msgLen == 0 {
				this.onMessage(this.msg)
			}
		}

		return nil
	}

	// if l > msgLen
	if this.msg != nil && this.onMessage != nil {
		this.msg.SetData(buf[:this.msgLen])
		this.onMessage(this.msg)
	}

	buf = buf[this.msgLen:]
	this.msgLen = 0

	return buf
}
