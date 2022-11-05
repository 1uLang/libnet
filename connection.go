package libnet

import (
	"crypto/tls"
	"errors"
	"github.com/1uLang/libnet/message"
	"github.com/1uLang/libnet/options"
	"github.com/1uLang/libnet/utils"
	"github.com/1uLang/libnet/utils/maps"
	"github.com/1uLang/libnet/workers"
	"github.com/mailru/easygo/netpoll"
	log "github.com/sirupsen/logrus"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	poller, _        = netpoll.New(nil)
	bytePool         = utils.NewBytePool(10_000, 65536)
	connId           = int64(0)
	countConnections = int64(0)
	connectionMaps   = map[int64]*Connection{}
	sharedLocker     = sync.Mutex{}
)

type Connection struct {
	conn       net.Conn
	desc       *netpoll.Desc
	worker     *workers.Worker
	buffer     *message.Buffer
	options    *options.Options
	isUdp      bool
	isClosed   bool
	remoteAddr string

	isClient bool
	connId   int64
	locker   sync.RWMutex
	context  maps.Map
	handler  Handler
}

func newConnection(rawConn net.Conn, handler Handler, opts *options.Options, isUdp, isClient bool) *Connection {
	conn := &Connection{
		isUdp:    isUdp,
		isClient: isClient,
		conn:     rawConn,
		options:  opts,
		handler:  handler,
		worker:   workers.Get(),
		context:  maps.Map{},
	}

	atomic.AddInt64(&countConnections, 1)
	sharedLocker.Lock()
	connId++
	conn.connId = connId
	connectionMaps[connId] = conn
	sharedLocker.Unlock()
	// 执行启动回调函数
	if !isUdp && conn.handler != nil && conn.handler.OnConnect != nil {
		conn.handler.OnConnect(conn)
	}
	return conn
}

// UDP 建立
func (this *Connection) setupUDP() {
	// 读取数据
	buf := bytePool.Get()
	defer bytePool.Put(buf)
	for {
		n, addr, err := this.conn.(*net.UDPConn).ReadFromUDP(buf)
		if err != nil {
			log.Error("[CONNECTION] read from error ", err)
			continue
		} else {
			this.remoteAddr = addr.String()
			if this.handler != nil && this.handler.OnConnect != nil {
				this.handler.OnConnect(this)
			}
		}
		if n > 0 {
			// udp client 不存在接受消息
			if !this.isClient && this.handler != nil && this.handler.OnMessage != nil {
				if this.options != nil && this.options.EncryptMethod != nil {
					decode, err := this.options.EncryptMethod.Decrypt(buf[:n])
					if err != nil {
						this.fail(errors.New("encryptMethod decrypt bytes fail:" + err.Error()))
					} else {
						this.handler.OnMessage(this, decode)
					}
				} else {
					this.handler.OnMessage(this, buf[:n])
				}
			}
		}
		// Close connection
		if this.handler != nil && this.handler.OnClose != nil {
			this.handler.OnClose(this, "")
		}
	}

}

// TCP 建立
func (this *Connection) setupTCP() {
	// 设置超时
	if this.options != nil && this.options.Timeout != 0 {
		this.conn.SetReadDeadline(time.Now().Add(this.options.Timeout))
	}
	this.remoteAddr = this.conn.RemoteAddr().String()
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
					_ = this.Close("timeout")
					return
				}
				if n > 0 {
					// 设置超时
					if this.options != nil && this.options.Timeout != 0 {
						this.conn.SetReadDeadline(time.Now().Add(this.options.Timeout))
					}
					if this.buffer != nil || this.handler != nil && this.handler.OnMessage != nil {
						if this.options != nil && this.options.EncryptMethod != nil {
							decode, err := this.options.EncryptMethod.Decrypt(buf[:n])
							if err != nil {
								this.fail(err)
							} else {
								if this.buffer != nil {
									this.buffer.Write(decode)
								} else {
									this.handler.OnMessage(this, decode)
								}
							}
						} else {
							if this.buffer != nil {
								this.buffer.Write(buf[:n])
							} else {
								this.handler.OnMessage(this, buf[:n])
							}
						}
					}
				} else {
					break
				}
			}
			bytePool.Put(buf)
			// 处理连接断开事件
			if ev&netpoll.EventReadHup != 0 {
				_ = this.Close("client hup")
				return
			}
		})
	})
	if err != nil {
		this.fail(err)
		return
	}
}

// TLS 建立
func (this *Connection) setupTLS() {
	// 设置超时
	if this.options != nil && this.options.Timeout != 0 {
		this.conn.SetReadDeadline(time.Now().Add(this.options.Timeout))
	}

	this.remoteAddr = this.conn.RemoteAddr().String()

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
				if err != nil && strings.Contains(err.Error(), "timeout") {
					bytePool.Put(buf)
					// 处理读取超时
					_ = this.Close("timeout")
					return
				}
				if n > 0 {
					// 设置超时
					if this.options != nil && this.options.Timeout != 0 {
						this.conn.SetReadDeadline(time.Now().Add(this.options.Timeout))
					}
					if this.buffer != nil || this.handler != nil && this.handler.OnMessage != nil {
						if this.options != nil && this.options.EncryptMethod != nil {
							decode, err := this.options.EncryptMethod.Decrypt(buf[:n])
							if err != nil {
								this.fail(err)
							} else {
								if this.buffer != nil {
									this.buffer.Write(decode)
								} else {
									this.handler.OnMessage(this, decode)
								}
							}
						} else {
							if this.buffer != nil {
								this.buffer.Write(buf[:n])
							} else {
								this.handler.OnMessage(this, buf[:n])
							}
						}
					}
				} else {
					break
				}
			}
			bytePool.Put(buf)
			// 处理连接断开事件
			if ev&netpoll.EventReadHup != 0 {
				_ = this.Close("client hup")
				return
			}
		})
	})
	if err != nil {
		this.fail(err)
		return
	}
}

// 异常
func (this *Connection) fail(err error) {
	log.Fatal(err)
}

// Close 主动断开连接
func (this *Connection) Close(reason string) error {
	defer func() {
		this.worker.Close()
	}()

	this.locker.Lock()
	if !this.isClosed {
		if this.handler != nil && this.handler.OnClose != nil {
			this.handler.OnClose(this, reason)
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

// IsClose 是否已断开
func (this *Connection) IsClose() bool {

	this.locker.RLocker()
	defer this.locker.RUnlock()
	return this.isClosed
}

// Write 下发消息
func (this *Connection) Write(bytes []byte) (n int, err error) {
	if this.IsClose() {
		return 0, errors.New("connection is close")
	}
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

// SetBuffer 设置接受消息监听器[注意当设置监听器之后 handler OnMessage将失效]
func (this *Connection) SetBuffer(buffer *message.Buffer) error {
	if this.IsClose() {
		return errors.New("the connection is close")
	}
	// udp client 不存在接受消息 股不存在设置接受消息监听器
	if this.isUdp && this.isClient {
		return errors.New("udp client is not to be set ")
	}
	this.buffer = buffer
	return nil
}

// RemoteAddr 远端地址
func (this *Connection) RemoteAddr() string {
	if this.remoteAddr == "" {
		return this.conn.RemoteAddr().String()
	}
	return this.remoteAddr
}
