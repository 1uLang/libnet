package example

import (
	"fmt"
	"github.com/1uLang/libnet/connection"
)

type Handle struct {
}

// OnConnect 当TCP长连接建立成功是回调
func (Handle) OnConnect(c *connection.Connection) {
	fmt.Println("new connection : ", c.RemoteAddr())
}

// OnMessage 当客户端有数据写入是回调
func (Handle) OnMessage(c *connection.Connection, bytes []byte) {
	fmt.Println("recv new msg : ", string(bytes))
	c.Write(bytes)
}

// OnClose 当客户端主动断开链接或者超时时回调,err返回关闭的原因
func (Handle) OnClose(c *connection.Connection, err error) {
	fmt.Println("close connection : ", c.RemoteAddr())
}
