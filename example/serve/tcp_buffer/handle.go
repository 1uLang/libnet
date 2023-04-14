package main

import (
	"fmt"
	"github.com/1uLang/libnet"
	message2 "github.com/1uLang/libnet/example/message"
	"github.com/1uLang/libnet/message"
)

type Handle struct{}
type conn struct {
	c      *libnet.Connection
	buffer *message.Buffer
}

func (c *conn) onMessage(m message.MessageI) {
	msg := m.(*message2.Message)
	msg.SetData([]byte("recv msg"))
	_, _ = c.c.Write(msg.Marshal())
}

// OnConnect 当TCP长连接建立成功是回调
func (Handle) OnConnect(c *libnet.Connection) {
	fmt.Println("new connection : ", c.RemoteAddr())
	buffer := message.NewBuffer(message2.CheckHeader)
	buffer.OnMessage((&conn{c: c}).onMessage)
	c.SetBuffer(buffer)
}

// OnMessage 当客户端有数据写入是回调
func (Handle) OnMessage(c *libnet.Connection, bytes []byte) {
	c.Write(bytes)
}

// OnClose 当客户端主动断开链接或者超时时回调,err返回关闭的原因
func (Handle) OnClose(c *libnet.Connection, msg string) {
	fmt.Println("close connection : ", c.RemoteAddr(), msg)
}
