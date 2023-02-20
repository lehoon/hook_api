package udp

import (
	"encoding/hex"
	"github.com/lehoon/hook_api/v2/library/logger"
	"net"
)

const (
	BindOk     = 1
	BindError  = 2
)

type IUdpServerEventCallback interface {
	OnBind()
	OnMessage(buffer []byte, length int) bool
	OnClose()
}

type DefaultUdpServer struct {
	addr     string
	bindFlag int8     //是否已经建立连接
	conn     *net.UDPConn //连接对象
	eventCallback IUdpServerEventCallback
}

func NewUdpServer(addr string, callback IUdpServerEventCallback) *DefaultUdpServer {
	return &DefaultUdpServer{
		addr:          addr,
		bindFlag:      BindError,
		conn:          nil,
		eventCallback: callback,
	}
}

//bind udp server
func (self *DefaultUdpServer) Bind() error {
	udpAddr, err := net.ResolveUDPAddr("udp", self.addr)
	if err != nil {
		logger.Log().Errorf("解析本地udp server地址失败")
		return err
	}

	self.conn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	self.bindFlag = BindOk

	if self.eventCallback != nil && self.eventCallback.OnBind != nil {
		self.eventCallback.OnBind()
	}

	go self.handleUdpConnection()
	return nil
}

func (self *DefaultUdpServer) handleUdpConnection()  {
	buffer := make([]byte, 2048)
	for {
		n, addr, err := self.conn.ReadFromUDP(buffer[0:])

		if err != nil {
			logger.Log().Error("读取数据失败.")
			return
		}

		if addr == nil {
			logger.Log().Error("读取数据失败.")
			continue
		}

		logger.Log().Infof("udp client Addr:[%s], data:[%s] count:[%d]", addr, hex.EncodeToString(buffer[:n]), n)
		if self.eventCallback != nil && self.eventCallback.OnMessage != nil {
			self.eventCallback.OnMessage(buffer, n)
		}
	}
}

func (self *DefaultUdpServer) Shutdown()  {
	self.conn.Close()
	self.bindFlag = BindError

	if self.eventCallback != nil && self.eventCallback.OnClose != nil {
		self.eventCallback.OnClose()
	}
}