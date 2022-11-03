package libnet

// 处理连接事件的回调接口

type Handler interface {
	OnConnect(c *Connection)               // 新连接回调
	OnMessage(c *Connection, bytes []byte) // 新消息回调
	OnClose(c *Connection, msg string)     // 连接断开回调
}
