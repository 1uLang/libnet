//go:build windows

/*
   @Author: 1usir
   @Description:
   @File: net_windows
   @Version: 1.0.0
   @Date: 2022/12/12 21:34
*/

package utils

import (
	"net"
	"syscall"
)

func ReadConn(syscallConn syscall.RawConn, buf []byte) (n int, err error) {
	return 0, nil
}

// 监听可重用的端口
func ListenReuseAddr(network string, addr string) (net.Listener, error) {
	return nil, nil
}

func SetLimit() {
}
