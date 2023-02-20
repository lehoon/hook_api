package net

import (
	"net"
	"time"
)

//检查远端服务是否可用
func CheckNetState(remoteHost, remotePort string) bool {
	address := net.JoinHostPort(remoteHost, remotePort)

	//retry connect
	conn, err := net.DialTimeout("tcp", address, time.Second * 3)
	if err != nil {
		return false
	}

	if conn == nil {
		return false
	}

	conn.Close()
	return true
}
