//go:build !windows && !android

package libnet

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/1uLang/libnet/message"
	"github.com/1uLang/libnet/options"
	"github.com/1uLang/libnet/utils"
	"github.com/1uLang/libnet/utils/maps"
	"github.com/1uLang/libnet/workers"
	"github.com/mailru/easygo/netpoll"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"path/filepath"
	"runtime"
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

	onClose func()
}

func newConnection(rawConn net.Conn, handler Handler, opts *options.Options, isUdp, isClient bool) *Connection {
	conn := &Connection{
		isUdp:    isUdp,
		isClient: isClient,
		conn:     rawConn,
		options:  opts,
		handler:  handler,
		worker:   workers.NewWorker(""),
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
	// conn
	desc, err := netpoll.Handle(this.conn, netpoll.EventRead|netpoll.EventEdgeTriggered)
	if err != nil {
		this.fail(errors.New("tls net poll handle " + err.Error()))
		_ = this.Close(err.Error())
		return
	}
	this.desc = desc

	syscallConn, err := this.conn.(*net.TCPConn).SyscallConn()
	if err != nil {
		this.fail(errors.New("tcp SyscallConn " + err.Error()))
		_ = this.Close(err.Error())
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
		this.fail(errors.New("tcp poller " + err.Error()))
		return
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

	tlsConn, _ := this.conn.(*tls.Conn)
	conn, _ := tlsConn.NetConn().(*net.TCPConn)
	// conn
	desc, err := netpoll.Handle(conn, netpoll.EventRead|netpoll.EventEdgeTriggered)
	if err != nil {
		this.fail(errors.New("tls net poll handle " + err.Error()))
		_ = this.Close(err.Error())
		return
	}
	this.desc = desc

	err = poller.Start(desc, func(ev netpoll.Event) {
		this.worker.Run(func() {
			// 读取数据
			buf := bytePool.Get()
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
			bytePool.Put(buf)
			// 处理连接断开事件
			if ev&netpoll.EventReadHup != 0 {
				_ = this.Close("client hup")
				return
			}
		})
	})
	if err != nil {
		this.fail(errors.New(this.remoteAddr + " tls poller " + err.Error()))
		return
	}
}

// 异常
func (this *Connection) fail(err error) {
	errorString := err.Error()

	// 调用stack
	_, currentFilename, _, currentOk := runtime.Caller(0)
	if currentOk {
		for i := 1; i < 32; i++ {
			_, filename, lineNo, ok := runtime.Caller(i)
			if !ok {
				break
			}

			if filename == currentFilename {
				continue
			}

			goPath := os.Getenv("GOPATH")
			if len(goPath) > 0 {
				absGoPath, err := filepath.Abs(goPath)
				if err == nil {
					filename = strings.TrimPrefix(filename, absGoPath)[1:]
				}
			} else if strings.Contains(filename, "src") {
				filename = filename[strings.Index(filename, "src"):]
			}

			errorString += "\n\t\t" + string(filename) + ":" + fmt.Sprintf("%d", lineNo)

			break
		}
	}

	log.Fatal(errorString)
}

// Close 主动断开连接
func (this *Connection) Close(reason string) error {

	if this.IsClose() {
		return nil
	}
	defer func() {
		if this.worker != nil {
			this.worker.Close()
		}
	}()
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
	// 关闭desc，需要在关闭conn之前
	if this.desc != nil {
		_ = poller.Stop(this.desc)
		_ = this.desc.Close()
	}
	return this.conn.Close()
}

// IsClose 是否已断开
func (this *Connection) IsClose() bool {
	this.locker.RLock()
	defer this.locker.RUnlock()
	return this.isClosed
}

// Write 下发消息
func (this *Connection) Write(bytes []byte) (n int, err error) {
	if this.IsClose() || this.conn == nil {
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
func (this *Connection) SetBuffer(buffer *message.Buffer) {
	// udp client 不存在接受消息 股不存在设置接受消息监听器
	if this.IsClose() || this.isUdp && this.isClient {
		return
	}
	this.buffer = buffer
	return
}

// 设置断开连接回调函数
func (this *Connection) SetOnClose(onClose func()) error {
	if this.IsClose() {
		return nil
	}
	this.onClose = onClose
	return nil
}

// RemoteAddr 远端地址
func (this *Connection) RemoteAddr() string {
	if this.remoteAddr == "" && this.conn != nil {
		return this.conn.RemoteAddr().String()
	}
	return this.remoteAddr
}
