package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/lehoon/hook_api/v2/library/logger"
)

// 录制 hls ts 切片完成事件
func RecordTsFinish(w http.ResponseWriter, r *http.Request) {
	result := SuccessBizResult()
	render.Respond(w, r, result)
	logger.Log().Info("record ts finish")
}

// 录制mp4切片完成事件
func RecordMP4Finish(w http.ResponseWriter, r *http.Request) {
	result := SuccessBizResult()
	render.Respond(w, r, result)
	logger.Log().Info("record mp4 finish")
}
