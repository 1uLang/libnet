//go:build windows || android

// windows 无epoll  故需要区分

package libnet

import (
	"errors"
	"github.com/1uLang/libnet/message"
	"github.com/1uLang/libnet/options"
	"github.com/1uLang/libnet/utils"
	"github.com/1uLang/libnet/utils/maps"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	bytePool         = utils.NewBytePool(10_000, 65536)
	connId           = int64(0)
	countConnections = int64(0)
	connectionMaps   = map[int64]*Connection{}
	sharedLocker     = sync.Mutex{}
)

type Connection struct {
	conn       net.Conn
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

	onClose func()
}

func newConnection(rawConn net.Conn, handler Handler, opts *options.Options, isUdp, isClient bool) *Connection {
	conn := &Connection{
		isUdp:    isUdp,
		isClient: isClient,
		conn:     rawConn,
		options:  opts,
		handler:  handler,
		context:  maps.Map{},
		locker:   sync.RWMutex{},
	}

	atomic.AddInt64(&countConnections, 1)
	sharedLocker.Lock()
	connId++
	conn.connId = connId
	connectionMaps[connId] = conn
	sharedLocker.Unlock()
	// 执行启动回调函数
	if !isUdp && conn.handler != nil {
		conn.handler.OnConnect(conn)
	}
	return conn
}

// UDP 建立
func (this *Connection) setupUDP() {
	// 读取数据
	buf := bytePool.Get()
	defer bytePool.Put(buf)

	if this.IsClose() {
		return
	}
	for {
		n, addr, err := this.conn.(*net.UDPConn).ReadFromUDP(buf)
		if err != nil {
			log.Error("[CONNECTION] read from error ", err)
			continue
		} else {
			this.remoteAddr = addr.String()
			if this.handler != nil {
				this.handler.OnConnect(this)
			}
		}
		if n > 0 {
			// udp client 不存在接受消息
			if !this.isClient && this.handler != nil {
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
		if this.handler != nil {
			this.handler.OnClose(this, "")
		}
	}
}

// TCP 建立
func (this *Connection) setupTCP() {

	if this.IsClose() {
		return
	}
	// 设置超时
	if this.options != nil && this.options.Timeout != 0 {
		this.conn.SetReadDeadline(time.Now().Add(this.options.Timeout))
	}
	this.remoteAddr = this.conn.RemoteAddr().String()
	// 读取数据
	buf := bytePool.Get()
	for {
		n, err := this.conn.Read(buf)
		if n > 0 {
			// 设置超时
			if this.options != nil && this.options.Timeout != 0 {
				this.conn.SetReadDeadline(time.Now().Add(this.options.Timeout))
			}
			if this.buffer != nil || this.handler != nil {
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
		}
		if err != nil {
			bytePool.Put(buf)
			if io.EOF == err {
				// 连接断开
				_ = this.Close("client close")
				return
			}
			if strings.Contains(err.Error(), "timeout") {
				// 读取超时
				_ = this.Close("timeout")
				return
			}
		}
	}
}

// TLS 建立
func (this *Connection) setupTLS() {

	if this.IsClose() {
		return
	}
	// 设置超时
	if this.options != nil && this.options.Timeout != 0 {
		this.conn.SetReadDeadline(time.Now().Add(this.options.Timeout))
	}

	this.remoteAddr = this.conn.RemoteAddr().String()
	// 读取数据
	buf := bytePool.Get()
	for {
		n, err := this.conn.Read(buf)
		if n > 0 {
			// 设置超时
			if this.options != nil && this.options.Timeout != 0 {
				this.conn.SetReadDeadline(time.Now().Add(this.options.Timeout))
			}
			if this.buffer != nil || this.handler != nil {
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
		}
		if err != nil {
			bytePool.Put(buf)
			if io.EOF == err {
				// 连接断开
				_ = this.Close("client close")
				return
			}
			if strings.Contains(err.Error(), "timeout") {
				// 读取超时
				_ = this.Close("timeout")
				return
			}
		}
	}
}

// 异常
func (this *Connection) fail(err error) {
	log.Fatal(err)
}

// Close 主动断开连接
func (this *Connection) Close(reason string) error {

	if this.IsClose() {
		return nil
	}

	this.locker.Lock()
	if !this.isClosed {
		if this.handler != nil {
			this.handler.OnClose(this, reason)
		}
		this.isClosed = true
		atomic.AddInt64(&countConnections, -1)
	}
	this.locker.Unlock()
	// 执行断开链接回调
	if this.onClose != nil {
		this.onClose()
	}
	err := this.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// IsClose 是否已断开
func (this *Connection) IsClose() bool {

	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.isClosed
}

// Write 下发消息
func (this *Connection) Write(bytes []byte) (n int, err error) {
	if this.IsClose() {
		return 0, nil
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
		return nil
	}
	// udp client 不存在接受消息 股不存在设置接受消息监听器
	if this.isUdp && this.isClient {
		return errors.New("udp client is not to be set ")
	}
	this.buffer = buffer
	return nil
}

// 设置断开连接回调函数
func (this *Connection) SetOnClose(onClose func()) error {
	if this.IsClose() {
		return errors.New("the connection is close")
	}
	this.onClose = onClose
	return nil
}

// RemoteAddr 远端地址
func (this *Connection) RemoteAddr() string {

	if this.remoteAddr == "" && !this.IsClose() {
		return this.conn.RemoteAddr().String()
	}
	return this.remoteAddr
}
