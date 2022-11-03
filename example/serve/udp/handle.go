package main

import (
	"fmt"
	"github.com/1uLang/libnet"
)

type Handle struct {
}

// OnConnect 当TCP长连接建立成功是回调
func (Handle) OnConnect(c *libnet.Connection) {
	fmt.Println("new connection : ", c.RemoteAddr())
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
