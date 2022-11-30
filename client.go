package libnet

import (
	"crypto/tls"
	"fmt"
	options2 "github.com/1uLang/libnet/options"
	"net"
)

type Client struct {
	options *options2.Options // 服务参数
	address string
	handler Handler
	conn    *Connection
}

func NewClient(address string, handler Handler, opts ...options2.Option) (*Client, error) {
	options := options2.GetOptions(opts...)
	setLimit()
	return &Client{
		options: options,
		address: address,
		handler: handler,
	}, nil
}
func (c *Client) Write(bytes []byte) (int, error) {
	if c.conn == nil {
		return 0, fmt.Errorf("not dial to server")
	}
	return c.conn.Write(bytes)
}
func (c *Client) DialTCP() error {
	rawConn, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}
	c.conn = newConnection(rawConn, c.handler, c.options, false, true)
	c.conn.setupTCP()
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
	if err != nil {
		return err
	}
	c.conn = newConnection(rawConn, c.handler, c.options, true, true)
	return nil
}
func (c *Client) DialTLS(cfg *tls.Config) error {
	var rawConn net.Conn
	var err error
	rawConn, err = tls.Dial("tcp", c.address, cfg)
	if err != nil {
		return err
	}
	c.conn = newConnection(rawConn, c.handler, c.options, false, true)
	c.conn.setupTLS()
	return nil
}
func (c *Client) Close() error {
	return c.conn.Close("")
}
