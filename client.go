package libnet

import (
	"crypto/tls"
	"fmt"
	"github.com/1uLang/libnet/connection"
	options2 "github.com/1uLang/libnet/options"
	"net"
)

type Client struct {
	options *options2.Options // 服务参数
	address string
	handler connection.Handler
	conn    *connection.Connection
}

func NewClient(address string, handler connection.Handler, opts ...options2.Option) (*Client, error) {
	options := options2.GetOptions(opts...)
	setLimit()
	if err := options2.CheckOptions(options); err != nil {
		return nil, fmt.Errorf("set options error : %s", err)
	}
	return &Client{
		options: options,
		address: address,
		handler: handler,
	}, nil
}
func (c *Client) Write(bytes []byte) (int, error) {
	if c.conn == nil {
		return 0, fmt.Errorf("not  dial to server")
	}
	return c.conn.Write(bytes)
}
func (c *Client) DialTCP() error {

	rawConn, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}
	c.conn = connection.NewConnection(rawConn, c.handler, c.options)
	// 执行启动回调函数
	if c.handler != nil && c.handler.OnConnect != nil {
		c.handler.OnConnect(c.conn)
	}
	c.conn.SetupTCP()
	return nil
}

func (c *Client) DialUDP() error {
	var rawConn net.Conn
	var err error
	if c.options != nil && c.options.Timeout != 0 {
		rawConn, err = net.DialTimeout("udp", c.address, c.options.Timeout)
	} else {
		rawConn, err = net.Dial("udp", c.address)
	}
	c.conn = connection.NewConnection(rawConn, c.handler, c.options)
	return err
}
func (c *Client) DialTLS(cfg *tls.Config) error {
	conn, err := tls.Dial("tcp", c.address, cfg)
	if err != nil {
		return err
	}
	c.conn = connection.NewConnection(conn, c.handler, c.options)
	// 执行启动回调函数
	if c.handler != nil && c.handler.OnConnect != nil {
		c.handler.OnConnect(c.conn)
	}
	c.conn.SetupTLS()
	return nil
}
