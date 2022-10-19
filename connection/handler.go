package connection

// Handler Server 注册接口
type Handler interface {
	OnConnect(c *Connection)               // OnConnect 当TCP长连接建立成功是回调
	OnMessage(c *Connection, bytes []byte) // OnMessage 当客户端有数据写入是回调
	OnClose(c *Connection, err error)      // OnClose 当客户端主动断开链接或者超时时回调,err返回关闭的原因
}
