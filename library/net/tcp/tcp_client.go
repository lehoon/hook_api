package tcp

import (
	"errors"
	"net"
	"time"

	"github.com/lehoon/hook_api/v2/library/logger"
)

const (
	ConnectClose = 1
	ConnectDoing = 2
	ConnectDone  = 3
)

type ITcpClientEventCallback interface {
	OnConnect()
	OnMessage(buffer []byte, length int) bool
	OnSend(length int) bool
	OnClose()
	OnDisConnect()
	OnReConnect()
}

type DefaultTcpClient struct {
	addr          string
	connFlag      int8     //是否已经建立连接
	conn          net.Conn //连接对象
	eventCallback ITcpClientEventCallback
}

func NewClient(addr string, callback ITcpClientEventCallback) *DefaultTcpClient {
	return &DefaultTcpClient{
		addr:          addr,
		connFlag:      ConnectClose,
		conn:          nil,
		eventCallback: callback}
}

func (client *DefaultTcpClient) IsConnect() bool {
	return client.connFlag == ConnectDone
}

func (client *DefaultTcpClient) IsConnecting() bool {
	return client.connFlag == ConnectDoing
}

// connect to server
func (client *DefaultTcpClient) Connect(timeout uint16) error {
	if client.addr == "" {
		return errors.New("服务器地址不能为空")
	}

	client.connFlag = ConnectDoing
	conn, err := net.DialTimeout("tcp", client.addr, time.Second*time.Duration(timeout))

	if err != nil {
		client.connFlag = ConnectClose
		logger.Log().Errorf("连接远程服务器%s失败%v", client.addr, err)
		return errors.New("连接远程服务器失败")
	}

	client.conn = conn
	client.connFlag = ConnectDone

	if client.eventCallback != nil && client.eventCallback.OnConnect != nil {
		client.eventCallback.OnConnect()
	}

	go client.handleConnection()
	return nil
}

func (client *DefaultTcpClient) handleConnection() {
	defer func() {
		err := client.conn.Close()
		if err != nil {
			logger.Log().Errorf("关闭socket连接失败 %v", err)
		}
	}()

	buffer := make([]byte, 4096)

	for {
		n, err := client.conn.Read(buffer)
		if err != nil {
			client.trigerOnDisConnect()
			return
		}

		if n > 0 {
			client.trigerOnMessage(buffer, n)
		}
	}
}

func (client *DefaultTcpClient) Send(buffer []byte, length int) (int, error) {
	if !client.IsConnect() {
		return 0, errors.New("当前与远程服务器未建立有效连接,发送数据失败")
	}

	idx := 0
	for {
		write, err := client.conn.Write(buffer[idx:])
		if err != nil {
			return idx, err
		}

		idx += write

		if idx >= length {
			break
		}
	}

	if client.eventCallback != nil && client.eventCallback.OnSend != nil {
		client.eventCallback.OnSend(idx)
	}

	return idx, nil
}

func (client *DefaultTcpClient) DisConnect() error {
	if !client.IsConnect() {
		return errors.New("当前与远程服务器未建立连接,关闭连接失败")
	}

	err := client.conn.Close()
	if err != nil {
		logger.Log().Errorf("关闭与远程服务器连接失败%v", err)
		return err
	}

	if client.eventCallback != nil && client.eventCallback.OnDisConnect != nil {
		client.eventCallback.OnDisConnect()
	}
	return nil
}

func (client *DefaultTcpClient) trigerOnMessage(buffer []byte, length int) {
	if client.eventCallback != nil && client.eventCallback.OnDisConnect != nil {
		client.eventCallback.OnMessage(buffer, length)
	}
}

func (client *DefaultTcpClient) trigerOnDisConnect() {
	if client.eventCallback != nil && client.eventCallback.OnDisConnect != nil {
		client.eventCallback.OnDisConnect()
	}

	client.connFlag = ConnectClose
}

func (client *DefaultTcpClient) trigerOnClose() {
	if client.eventCallback != nil && client.eventCallback.OnClose != nil {
		client.eventCallback.OnClose()
	}
}

// 远程服务是否在线
func IsOnline(address string) error {
	conn, err := net.DialTimeout("tcp", address, time.Second*time.Duration(3000))

	if err != nil {
		logger.Log().Errorf("连接远程服务器%s失败%v", address, err)
		return errors.New("连接远程服务器失败")
	}

	conn.Close()
	return nil
}
