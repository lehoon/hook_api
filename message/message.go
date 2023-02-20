package message

import (
	"errors"
	"net/http"
	"strings"
)

type Message struct {
	MediaServiceId string `json:"mediaServiceId"`
}

// 流不存在请求
type StreamNotFoundMessage struct {
	Id     string
	App    string
	Ip     string
	Port   int
	Params string
	Stream string
	Vhost  string
}

type StreamNotFoundBind struct {
	*StreamNotFoundMessage
}

func (a *StreamNotFoundBind) Bind(r *http.Request) error {
	if a.StreamNotFoundMessage == nil {
		return errors.New("没有找到请求参数")
	}

	return nil
}

// 心跳请求
type KeepAliveReportMessage struct {
	Buffer                uint64
	BufferLikeString      uint64
	BufferList            uint64
	BufferRaw             uint64
	Frame                 uint64
	FrameImp              uint64
	MediaSource           uint64
	MultiMediaSourceMuxer uint64
	RtmpPacket            uint64
	RtpPacket             uint64
	Socket                uint64
	TcpClient             uint64
	TcpServer             uint64
	TcpSession            uint64
	UdpServer             uint64
	UdpSession            uint64
}

type KeepAliveReportRequest struct {
	MediaServiceId string `json:"mediaServiceId"`
	Data           KeepAliveReportMessage
}

type KeepAliveReportBind struct {
	*KeepAliveReportRequest
}

func (a *KeepAliveReportBind) Bind(r *http.Request) error {
	if a.KeepAliveReportRequest == nil {
		return errors.New("没有找到请求参数")
	}

	return nil
}

// 设备信息
type DeviceInfo struct {
	StreamId  string `json:"streamId"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Hostname  string `json:"hostname"`
	AppName   string `json:"appName"`
	VHostName string `json:"vhostName"`
}

func (d *DeviceInfo) IsEmpty() bool {
	return len(d.AppName) == 0 && len(d.Hostname) == 0 && len(d.VHostName) == 0 && len(d.Password) == 0
}

func (d *DeviceInfo) CanPublish() bool {
	return len(d.AppName) > 0 && len(d.Hostname) > 0 && len(d.VHostName) > 0 && len(d.Password) > 0 && len(d.StreamId) > 0
}

func (d *DeviceInfo) PullStreamKey() string {
	var builder strings.Builder
	builder.WriteString("&vhost=")
	builder.WriteString(d.VHostName)
	builder.WriteString("&app=")
	builder.WriteString(d.AppName)
	builder.WriteString("&stream=")
	builder.WriteString(d.StreamId)
	builder.WriteString("&url=rtsp://")
	builder.WriteString(d.Username)
	builder.WriteString(":")
	builder.WriteString(d.Password)
	builder.WriteString("@")
	builder.WriteString(d.Hostname)
	return builder.String()
}

func (d *DeviceInfo) IsOnlineKey() string {
	var builder strings.Builder
	builder.WriteString("vhost=")
	builder.WriteString(d.VHostName)
	builder.WriteString("&app=")
	builder.WriteString(d.AppName)
	builder.WriteString("&stream=")
	builder.WriteString(d.StreamId)
	return builder.String()
}

func (d *DeviceInfo) CloseKey() string {
	var builder strings.Builder
	builder.WriteString(d.VHostName)
	builder.WriteString("/")
	builder.WriteString(d.AppName)
	builder.WriteString("/")
	builder.WriteString(d.StreamId)
	return builder.String()
}

type DeviceInfoBind struct {
	*DeviceInfo
}

func (a *DeviceInfoBind) Bind(r *http.Request) error {
	if a.DeviceInfo == nil {
		return errors.New("没有找到请求参数")
	}

	return nil
}

// 检查流是否在线响应消息
type StreamIsOnlineResponse struct {
	Code   int  `json:"code"`
	Online bool `json:"online"`
}
