package message

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/lehoon/hook_api/v2/library/config"
)

const (
	STREAM_TRANSPORT_PROTOCOL_RTSP = 0 //传输协议  rtsp
)

const (
	STREAM_TRANSPORT_PROTOCOL_RTSP_PORT_DEFAULT = 554 //RTSP port
)

type Message struct {
	MediaServiceId string `json:"mediaServiceId"`
}

// 流不存在请求
type StreamNotFoundMessage struct {
	Id       string `json:"id"`
	App      string `json:"app"`
	Ip       string `json:"ip"`
	Port     int    `json:"port"`
	Params   string `json:"params"`
	StreamId string `json:"stream"`
	Vhost    string `json:"vhost"`
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

func (snfb *StreamNotFoundBind) JsonString() string {
	buf, err := json.Marshal(snfb)
	if err != nil {
		return ""
	}

	return string(buf)
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
	DeviceId  string `json:"deviceId"`
	StreamId  string `json:"-"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Hostname  string `json:"hostname"`
	AppName   string `json:"appName"`
	VHostName string `json:"vhostName"`
	Protocol  uint8  `json:"protocol"`
	Port      uint16 `json:"port"`
}

func (d *DeviceInfo) IsEmpty() bool {
	return len(d.AppName) == 0 && len(d.Hostname) == 0 && len(d.VHostName) == 0 && len(d.Password) == 0
}

func (d *DeviceInfo) CanPublish() bool {
	return len(d.AppName) > 0 && len(d.Hostname) > 0 && len(d.VHostName) > 0 && len(d.Password) > 0 && len(d.DeviceId) > 0
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

// 流信息
type StreamInfo struct {
	DeviceId  string `json:"-"`
	StreamId  string `json:"streamId"`
	Username  string `json:"-"` //`json:"username"`
	Password  string `json:"-"` //`json:"password"`
	Hostname  string `json:"-"` //`json:"hostname"`
	AppName   string `json:"appName"`
	VHostName string `json:"-"` //`json:"vhostName"`
	Protocol  uint8  `json:"protocol"`
	Port      uint16 `json:"port"`
}

func (s *StreamInfo) stream_play_url() string {
	var builder strings.Builder
	if s.Protocol == STREAM_TRANSPORT_PROTOCOL_RTSP {
		builder.WriteString("rtsp://")
		builder.WriteString(s.Username)
		builder.WriteString(":")
		builder.WriteString(s.Password)
		builder.WriteString("@")
		builder.WriteString(s.Hostname)
		builder.WriteString(":")
		builder.WriteString(strconv.FormatInt(int64(s.Port), 10))
	}

	return builder.String()
}

func (d *StreamInfo) PullStreamKey() string {
	var builder strings.Builder
	builder.WriteString("&vhost=")
	builder.WriteString(d.VHostName)
	builder.WriteString("&app=")
	builder.WriteString(d.AppName)
	builder.WriteString("&stream=")
	builder.WriteString(d.StreamId)
	builder.WriteString("&url=")
	builder.WriteString(d.stream_play_url())
	//builder.WriteString("&url=rtsp://")
	//builder.WriteString(d.Username)
	//builder.WriteString(":")
	//builder.WriteString(d.Password)
	//builder.WriteString("@")
	//builder.WriteString(d.Hostname)
	builder.WriteString("&enable_rtmp=1")
	builder.WriteString("&enable_audio=0")
	return builder.String()
}

func (d *StreamInfo) IsOnlineKey() string {
	var builder strings.Builder
	builder.WriteString("vhost=")
	builder.WriteString(d.VHostName)
	builder.WriteString("&app=")
	builder.WriteString(d.AppName)
	builder.WriteString("&stream=")
	builder.WriteString(d.StreamId)
	return builder.String()
}

func (d *StreamInfo) CloseKey() string {
	var builder strings.Builder
	builder.WriteString(d.VHostName)
	builder.WriteString("/")
	builder.WriteString(d.AppName)
	builder.WriteString("/")
	builder.WriteString(d.StreamId)
	return builder.String()
}

func (d *StreamInfo) PlayUrl() []string {
	var result []string
	result = append(result, "rtsp://"+config.GetServerPlayUrl()+"/"+d.AppName+"/"+d.StreamId)
	result = append(result, "http://"+config.GetServerPlayUrl()+"/"+d.AppName+"/"+d.StreamId+".live.flv")
	result = append(result, "ws://"+config.GetServerPlayUrl()+"/"+d.AppName+"/"+d.StreamId+".live.flv")
	return result
}

// 检查流是否在线响应消息
type StreamIsOnlineResponse struct {
	Code   int  `json:"code"`
	Online bool `json:"online"`
}

// 根据协议获取端口号
func GetProtocolPort(protocol uint8, port uint16) uint16 {
	if port > 0 && port < 65535 {
		return port
	}

	if protocol == STREAM_TRANSPORT_PROTOCOL_RTSP {
		return STREAM_TRANSPORT_PROTOCOL_RTSP_PORT_DEFAULT
	}

	return 0
}
