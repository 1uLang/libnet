package main

import (
	"fmt"
	"github.com/1uLang/libnet"
	"github.com/1uLang/libnet/message"
)

type Handle struct{}
type conn struct {
	c *libnet.Connection
}

func (c *conn) onMessage(msg *message.Message) {
	fmt.Println("recv msg : ", *msg)
	msg.Data = []byte("recv msg")
	n, err := c.c.Write(msg.Marshal())
	fmt.Println("send msg ", n, err)
}

// OnConnect 当TCP长连接建立成功是回调
func (Handle) OnConnect(c *libnet.Connection) {
	fmt.Println("new connection : ", c.RemoteAddr())
	buffer := message.NewBuffer()
	buffer.OnMessage((&conn{c: c}).onMessage)
	c.SetBuffer(buffer)
}

// OnMessage 当客户端有数据写入是回调
func (Handle) OnMessage(c *libnet.Connection, bytes []byte) {
	fmt.Println("recv new msg : ", string(bytes))
	c.Write(bytes)
}

// OnClose 当客户端主动断开链接或者超时时回调,err返回关闭的原因
func (Handle) OnClose(c *libnet.Connection, msg string) {
	fmt.Println("close connection : ", c.RemoteAddr(), msg)
}
