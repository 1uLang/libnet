//go:build !linux
// +build !linux

// 可以在 /usr/include/asm-generic/socket.h 中找到 SO_REUSEPORT 值

package utils

import (
	"net"
	"syscall"
)

const SO_REUSEPORT = 15

// 从连接中读取数据
func ReadConn(syscallConn syscall.RawConn, buf []byte) (n int, err error) {
	return 0, nil
}

// 监听可重用的端口
func ListenReuseAddr(network string, addr string) (net.Listener, error) {
	return nil, nil
}
