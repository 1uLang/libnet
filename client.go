package libnet

import (
	"crypto/tls"
	"fmt"
	options2 "github.com/1uLang/libnet/options"
	"github.com/1uLang/libnet/utils"
	"net"
	"time"
)

type Client struct {
	options *options2.Options // 服务参数
	address string
	handler Handler
	conn    *Connection
}

func NewClient(address string, handler Handler, opts ...options2.Option) (*Client, error) {
	options := options2.GetOptions(opts...)
	utils.SetLimit()
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

	timeout := c.options.Timeout
	if c.options.Timeout == 0 { // 默认3秒
		timeout = 3 * time.Second
	}
	rawConn, err := net.DialTimeout("tcp", c.address, timeout)
	if err != nil {
		return err
	}
	//rawConn.(*net.TCPConn).SetLinger(0)
	c.conn = newConnection(rawConn, c.handler, c.options, false, true)
	c.conn.setupTCP()
	return nil
}

func (c *Client) DialUDP() error {

	timeout := c.options.Timeout
	if c.options.Timeout == 0 { // 默认3秒
		timeout = 3 * time.Second
	}
	rawConn, err := net.DialTimeout("udp", c.address, timeout)
	if err != nil {
		return err
	}
	//err = rawConn.(*net.UDPConn).SetWriteBuffer(4 * 1024 * 1024)
	c.conn = newConnection(rawConn, c.handler, c.options, true, true)
	return nil
}
func (c *Client) DialTLS(cfg *tls.Config) error {
	var rawConn net.Conn
	var err error
	timeout := c.options.Timeout
	if c.options.Timeout == 0 { // 默认3秒
		timeout = 3 * time.Second
	}
	rawConn, err = tls.DialWithDialer(&net.Dialer{Timeout: timeout}, "tcp", c.address, cfg)
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
