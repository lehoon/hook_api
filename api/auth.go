package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/lehoon/hook_api/v2/library/logger"
)

// 访问http文件鉴权事件
func AuthHttpFile(w http.ResponseWriter, r *http.Request) {
	result := SuccessBizResult()
	render.Respond(w, r, result)
	logger.Log().Info("访问http文件鉴权事件")
}

// 播放鉴权
func AuthPlay(w http.ResponseWriter, r *http.Request) {
	result := SuccessBizResult()
	render.Respond(w, r, result)
	logger.Log().Info("播放鉴权")
}

// 推流鉴权事件
func AuthPublish(w http.ResponseWriter, r *http.Request) {
	result := SuccessBizResult()
	render.Respond(w, r, result)
	logger.Log().Info("推流鉴权事件")
}

// rtsp播放鉴权事件
func AuthRtspPlay(w http.ResponseWriter, r *http.Request) {
	result := SuccessBizResult()
	render.Respond(w, r, result)
	logger.Log().Info("rtsp播放鉴权事件")
}

// 远程telnet调试鉴权事件
func AuthShell(w http.ResponseWriter, r *http.Request) {
	result := SuccessBizResult()
	render.Respond(w, r, result)
	logger.Log().Info("远程telnet调试鉴权事件")
}

// rtsp播放是否开启专属鉴权事件
func IsRtspAuth(w http.ResponseWriter, r *http.Request) {
	result := SuccessBizResult()
	render.Respond(w, r, result)
	logger.Log().Info("rtsp播放是否开启专属鉴权事件")
}
