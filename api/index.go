package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/lehoon/hook_api/v2/library/logger"
)

func QueryOperateCodeMessage(w http.ResponseWriter, r *http.Request) {
	data := OperateCodeMessage()
	render.Respond(w, r, SuccessBizResultWithData(data))
}

func ShowPostMessage(w http.ResponseWriter, r *http.Request) {
	logger.Log().Infof("ShowPostMessage   %s", r.Body)
	data := OperateCodeMessage()
	render.Respond(w, r, SuccessBizResultWithData(data))
}
