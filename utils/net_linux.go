//go:build linux
// +build linux

// 可以在 /usr/include/asm-generic/socket.h 中找到 SO_REUSEPORT 值

package utils

import (
	"context"
	"fmt"
	"net"
	"syscall"
)

const SO_REUSEPORT = 15

// 从连接中读取数据
func ReadConn(syscallConn syscall.RawConn, buf []byte) (n int, err error) {
	err = syscallConn.Read(func(fd uintptr) (done bool) {
		err = syscall.SetNonblock(int(fd), true)
		if err != nil {
			n = -1
			return true
		}
		n, err = syscall.Read(int(fd), buf)
		if err != nil {
			n = -1
			return true
		}

		return true
	})
	return
}

// 监听可重用的端口
func ListenReuseAddr(network string, addr string) (net.Listener, error) {
	config := &net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, SO_REUSEPORT, 1)
				if err != nil {
					fmt.Println(err)
				}
			})
		},
		KeepAlive: 0,
	}
	return config.Listen(context.Background(), network, addr)
}
