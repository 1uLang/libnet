package connection

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/1uLang/libnet/options"
	"github.com/1uLang/libnet/utils"
	"github.com/1uLang/libnet/workers"
	"github.com/mailru/easygo/netpoll"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	poller, _        = netpoll.New(nil)
	bytePool         = utils.NewBytePool(10_000, 65536)
	countConnections int64
)

type Connection struct {
	worker     *workers.Worker
	desc       *netpoll.Desc
	conn       net.Conn
	handler    Handler
	options    *options.Options
	isClosed   bool
	remoteAddr string
	locker     sync.Mutex
}

func NewConnection(rawConn net.Conn, handler Handler, options *options.Options) *Connection {

	atomic.AddInt64(&countConnections, 1)

	return &Connection{
		worker:  workers.NewWorker(""),
		conn:    rawConn,
		handler: handler,
		options: options,
	}
}

func (this *Connection) SetupUDP() {
	// 读取数据
	var buf = make([]byte, 1024)
	for {
		n, addr, err := this.conn.(*net.UDPConn).ReadFromUDP(buf)
		if err != nil {
			utils.Log().Error("[CONNECTION] read from error ", err)
		} else {
			this.locker.Lock()
			this.remoteAddr = addr.String()
			this.locker.Unlock()
			if this.handler != nil && this.handler.OnConnect != nil {
				this.handler.OnConnect(this)
			}
		}
		if n > 0 {
			if this.handler != nil && this.handler.OnMessage != nil {
				if this.options != nil && this.options.EncryptMethod != nil {
					decode, err := this.options.EncryptMethod.Decrypt(buf[:n])
					if err != nil {
						this.fail(err)
					} else {
						this.handler.OnMessage(this, decode)
					}
				} else {
					this.handler.OnMessage(this, buf[:n])
				}
			}
		}
		// close connection
		if this.handler != nil && this.handler.OnClose != nil {
			this.handler.OnClose(this, nil)
		}
	}
}

func (this *Connection) SetupTCP() {

	this.locker.Lock()
	this.remoteAddr = this.conn.RemoteAddr().String()
	this.locker.Unlock()
	// conn
	desc, err := netpoll.Handle(this.conn, netpoll.EventRead|netpoll.EventEdgeTriggered)
	if err != nil {
		this.fail(err)
		return
	}
	this.desc = desc

	syscallConn, err := this.conn.(*net.TCPConn).SyscallConn()
	if err != nil {
		this.fail(err)
		return
	}

	err = poller.Start(desc, func(ev netpoll.Event) {
		this.worker.Run(func() {
			// 读取数据
			buf := bytePool.Get()
			for {
				n, err := utils.ReadConn(syscallConn, buf)
				if err != nil && strings.Contains(err.Error(), "timeout") {
					bytePool.Put(buf)
					// 处理读取超时
					_ = this.close("timeout")
					return
				}
				if n > 0 {
					// 设置超时
					if this.options != nil && this.options.Timeout != 0 {
						this.conn.SetReadDeadline(time.Now().Add(this.options.Timeout))
					}
					if this.handler != nil && this.handler.OnMessage != nil {
						if this.options != nil && this.options.EncryptMethod != nil {
							decode, err := this.options.EncryptMethod.Decrypt(buf[:n])
							if err != nil {
								this.fail(err)
							} else {
								this.handler.OnMessage(this, decode)
							}
						} else {
							this.handler.OnMessage(this, buf[:n])
						}
					}
				} else {
					break
				}
			}
			bytePool.Put(buf)
			// 处理连接断开事件
			if ev&netpoll.EventReadHup != 0 {
				_ = this.close("client hup")
				return
			}
		})
	})
	if err != nil {
		this.fail(err)
		return
	}
}
func (this *Connection) SetupTLS() {

	this.locker.Lock()
	this.remoteAddr = this.conn.RemoteAddr().String()
	this.locker.Unlock()

	tlsConn, _ := this.conn.(*tls.Conn)
	conn, _ := tlsConn.NetConn().(*net.TCPConn)
	// conn
	desc, err := netpoll.Handle(conn, netpoll.EventRead|netpoll.EventEdgeTriggered)
	if err != nil {
		this.fail(err)
		return
	}
	this.desc = desc

	err = poller.Start(desc, func(ev netpoll.Event) {
		this.worker.Run(func() {
			// 读取数据
			buf := bytePool.Get()
			for {
				n, err := tlsConn.Read(buf)
				fmt.Println("====== ", conn.RemoteAddr(), n, err, ev)
				if err != nil && strings.Contains(err.Error(), "timeout") {
					bytePool.Put(buf)
					// 处理读取超时
					_ = this.close("timeout")
					return
				}
				if n > 0 {
					// 设置超时
					if this.options != nil && this.options.Timeout != 0 {
						this.conn.SetReadDeadline(time.Now().Add(this.options.Timeout))
					}
					if this.handler != nil && this.handler.OnMessage != nil {
						if this.options != nil && this.options.EncryptMethod != nil {
							decode, err := this.options.EncryptMethod.Decrypt(buf[:n])
							if err != nil {
								this.fail(err)
							} else {
								this.handler.OnMessage(this, decode)
							}
						} else {
							this.handler.OnMessage(this, buf[:n])
						}
					}
				} else {
					break
				}
			}
			bytePool.Put(buf)
			// 处理连接断开事件
			if ev&netpoll.EventReadHup != 0 {
				_ = this.close("client hup")
				return
			}

		})
	})
	if err != nil {
		this.fail(err)
		return
	}
}

func (this *Connection) close(reason string) error {
	defer func() {
		this.worker.Close()
	}()

	this.locker.Lock()
	if !this.isClosed {
		if this.handler != nil && this.handler.OnClose != nil {
			this.handler.OnClose(this, errors.New(reason))
		}
		this.isClosed = true
		atomic.AddInt64(&countConnections, -1)
	}
	this.locker.Unlock()
	// 关闭desc，需要在关闭conn之前
	if this.desc != nil {
		_ = poller.Stop(this.desc)
		_ = this.desc.Close()
	}
	err := this.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (this *Connection) Write(bytes []byte) (n int, err error) {

	if this.options != nil && this.options.EncryptMethod != nil {
		bytes, err := this.options.EncryptMethod.Encrypt(bytes)
		if err != nil {
			this.fail(err)
		}
		return this.conn.Write(bytes)
	} else {
		return this.conn.Write(bytes)
	}
}

// 失败提示—
func (this *Connection) fail(err error) {
	utils.Log().Error("[CONNECTION] ", err)
}

func (this *Connection) RemoteAddr() string {
	this.locker.Lock()
	defer this.locker.Unlock()
	if this.remoteAddr == "" {
		return this.conn.RemoteAddr().String()
	}
	return this.remoteAddr
}
