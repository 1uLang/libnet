package libnet

import (
	"crypto/tls"
	"github.com/1uLang/libnet/options"
	"github.com/1uLang/libnet/utils"
	log "github.com/sirupsen/logrus"
	"net"
)

type Serve struct {
	address string
	// 服务参数
	options *options.Options
	// 处理消息回调接口
	handler Handler
}

func NewServe(address string, handler Handler, opts ...options.Option) *Serve {
	utils.SetLimit()
	return &Serve{
		options: options.GetOptions(opts...),
		address: address,
		handler: handler,
	}
}

func (s *Serve) RunUDP() error {
	log.Info("[Serve] Run ", s.address, " udp server")
	udpAddr, err := net.ResolveUDPAddr("udp", s.address)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		return err
	}

	newConnection(conn, s.handler, s.options, true, false).setupUDP()
	return nil
}

func (s *Serve) RunTCP() error {
	log.Info("[Serve] Run ", s.address, " tcp server")
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error("new tcp client connection error ", err)
		}
		newConnection(conn, s.handler, s.options, false, false).setupTCP()
	}
}

func (s *Serve) RunTLS(cfg *tls.Config) error {

	log.Info("[Serve] Run ", s.address, " tls server")

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
			log.Error("new tls client connection error ", err)
		}
		newConnection(conn.(*tls.Conn), s.handler, s.options, false, false).setupTLS()
	}
}
