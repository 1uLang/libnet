package libnet

import (
	"crypto/tls"
	"fmt"
	"github.com/1uLang/libnet/connection"
	options2 "github.com/1uLang/libnet/options"
	"github.com/1uLang/libnet/utils"
	"net"
	"syscall"
	"time"
)

type Serve struct {
	options *options2.Options // 服务参数
	address string
	handler connection.Handler
}

func NewServe(address string, handler connection.Handler, opts ...options2.Option) (*Serve, error) {
	options := options2.GetOptions(opts...)
	setLimit()
	if err := options2.CheckOptions(options); err != nil {
		return nil, fmt.Errorf("set options error : %s", err)
	}
	return &Serve{
		options: options,
		address: address,
		handler: handler,
	}, nil
}

func (s *Serve) RunTCP() error {
	utils.Log().Info("[Serve] Run ", s.address, "Tcp Server")
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			utils.Log().Error("new tcp client connection error ", err)
		}
		// 设置超时
		if s.options != nil && s.options.Timeout != 0 {
			conn.SetReadDeadline(time.Now().Add(s.options.Timeout))
		}
		c := connection.NewConnection(conn, s.handler, s.options)
		// 执行启动回调函数
		if s.handler != nil && s.handler.OnConnect != nil {
			s.handler.OnConnect(c)
		}
		c.SetupTCP()
	}
}
func (s *Serve) RunUDP() error {
	utils.Log().Info("[Serve] Run ", s.address, "Udp Server")
	udpAddr, err := net.ResolveUDPAddr("udp", s.address)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	if err != nil {
		return err
	}

	c := connection.NewConnection(conn, s.handler, s.options)

	c.SetupUDP()
	return nil
}
func (s *Serve) RunTLS(cfg *tls.Config) error {

	utils.Log().Info("[Serve] Run ", s.address, "TLS Server")

	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	defer ln.Close()
	tlsListener := tls.NewListener(ln, cfg)
	defer tlsListener.Close()
	for {
		conn, err := tlsListener.Accept()
		if err != nil {
			utils.Log().Error("new tcp client connection error ", err)
		}
		// 设置超时
		if s.options != nil && s.options.Timeout != 0 {
			conn.SetReadDeadline(time.Now().Add(s.options.Timeout))
		}
		c := connection.NewConnection(conn.(*tls.Conn), s.handler, s.options)
		// 执行启动回调函数
		if s.handler != nil && s.handler.OnConnect != nil {
			s.handler.OnConnect(c)
		}
		c.SetupTLS()
	}
}
func setLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

}
