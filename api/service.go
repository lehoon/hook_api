package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/lehoon/hook_api/v2/library/logger"
	"github.com/lehoon/hook_api/v2/message"
)

// 服务器启动报告，可以用于服务器的崩溃重启事件监听
func ServerStartupReport(w http.ResponseWriter, r *http.Request) {
	//buf, _ := io.ReadAll(r.Body)
	//logger.Log().Infof("服务器启动报告:%s", string(buf))
	logger.Log().Info("服务器启动报告")
	render.Respond(w, r, SuccessBizResult())
}

// server保活上报
func KeepAliveReport(w http.ResponseWriter, r *http.Request) {
	request := &message.KeepAliveReportBind{}
	if err := render.Bind(r, request); err != nil {
		logger.Log().Errorf("Service心跳请求失败, 获取请求参数失败, %s", err.Error())
		render.Render(w, r, FailureBizResultWithParamError())
		return
	}

	logger.Log().Info("server保活上报")
	//render.Render(w, r, SuccessBizResult())
	render.Respond(w, r, SuccessBizResult())
}

// 播放器或推流器使用流量事件
func FlowReport(w http.ResponseWriter, r *http.Request) {
	logger.Log().Info("播放器或推流器使用流量事件")
	render.Respond(w, r, SuccessBizResult())
}

// 发送rtp(startSendRtp)被动关闭时回调
func RtpCloseReport(w http.ResponseWriter, r *http.Request) {
	logger.Log().Info("发送rtp(startSendRtp)被动关闭时回调")
	render.Respond(w, r, SuccessBizResult())
}

// rtp server 超时未收到数据
func RtpTimeoutReport(w http.ResponseWriter, r *http.Request) {
	logger.Log().Info("rtp server 超时未收到数据")
	render.Respond(w, r, SuccessBizResult())
}
