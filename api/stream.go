package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lehoon/hook_api/v2/library/config"
	"github.com/lehoon/hook_api/v2/library/logger"
	restapi "github.com/lehoon/hook_api/v2/library/net/http"
	"github.com/lehoon/hook_api/v2/library/net/tcp"
	"github.com/lehoon/hook_api/v2/message"
	"github.com/lehoon/hook_api/v2/service"
)

func StreamList(w http.ResponseWriter, r *http.Request) {
	resultList, err := service.QueryStreamList()
	if err != nil {
		logger.Log().Errorf("查询流列表失败, %s", err.Error())
		render.Render(w, r, FailureBizResultWithDatabaseError())
		return
	}

	logger.Log().Info("查询流列表成功")
	render.Respond(w, r, SuccessBizResultWithData(resultList))
}

// 检查流是否可以播放
func StreamCanPlay(w http.ResponseWriter, r *http.Request) {

}

// 检查流是否存在
func StreamIsOnline(w http.ResponseWriter, r *http.Request) {
	streamId := chi.URLParam(r, "streamId")
	if streamId == "" {
		render.Respond(w, r, FailureBizResultWithParamError())
		return
	}

	if err := tcp.IsOnline(config.GetRestAddress()); err != nil {
		logger.Log().Errorf("查询流是否在线失败, %s, %s", streamId, err.Error())
		render.Render(w, r, FailureBizResultWithServiceNotOnline())
		return
	}

	isOnline, err := invokeStreamIsOnline(streamId)
	if err != nil {
		render.Respond(w, r, FailureBizResultWithMessage(OPEARTE_STREAM_ISONLINE_ERROR, err.Error()))
		return
	}

	logger.Log().Infof("查询流是否在线完成, %s,%v", streamId, isOnline)
	render.Respond(w, r, SuccessBizResultWithData(isOnline))
}

// 查询流的播放地址
func StreamPlayUrl(w http.ResponseWriter, r *http.Request) {
	streamId := chi.URLParam(r, "streamId")
	if streamId == "" {
		render.Respond(w, r, FailureBizResultWithParamError())
		return
	}

	//校验流编号在数据库是否已存在
	streamInfo, err := service.QueryStreamByStreamId(streamId)
	if err != nil || streamInfo == nil {
		logger.Log().Errorf("查询流播放地址失败, %s, %s", streamId, err.Error())
		render.Respond(w, r, FailureBizResultWithStreamNotExists())
		return
	}

	logger.Log().Infof("查询流播放地址成功, %s", streamId)
	render.Respond(w, r, SuccessBizResultWithData(streamInfo.PlayUrl()))
}

// 打开指定流
func StreamOpen(w http.ResponseWriter, r *http.Request) {
	streamId := chi.URLParam(r, "streamId")
	if streamId == "" {
		render.Respond(w, r, FailureBizResultWithParamError())
		return
	}

	if err := tcp.IsOnline(config.GetRestAddress()); err != nil {
		logger.Log().Errorf("打开流失败, %s, %s", streamId, err.Error())
		render.Render(w, r, FailureBizResultWithServiceNotOnline())
		return
	}

	//校验流编号在数据库是否已存在
	streamInfo, err := service.QueryStreamByStreamId(streamId)
	if err != nil || streamInfo == nil {
		logger.Log().Errorf("打开流失败, %s, %s", streamId, err.Error())
		render.Respond(w, r, FailureBizResultWithStreamNotExists())
		return
	}

	request := &message.StreamNotFoundBind{
		StreamNotFoundMessage: &message.StreamNotFoundMessage{
			App:    "",
			Vhost:  "",
			Stream: streamId,
		},
	}

	err = invokePullStream(request)

	if err != nil {
		logger.Log().Errorf("打开流成功失败, %s, %s", request.Stream, err.Error())
		render.Render(w, r, FailureBizResultWithMessage(OPEARTE_STREAM_NOT_EXITS, err.Error()))
		return
	}

	logger.Log().Infof("打开流成功, %s", streamId)
	render.Respond(w, r, SuccessBizResult())
}

// 关闭指定流
func StreamClose(w http.ResponseWriter, r *http.Request) {
	streamId := chi.URLParam(r, "streamId")
	if streamId == "" {
		render.Respond(w, r, FailureBizResultWithParamError())
		return
	}

	if err := tcp.IsOnline(config.GetRestAddress()); err != nil {
		logger.Log().Errorf("关闭流失败, %s, %s", streamId, err.Error())
		render.Render(w, r, FailureBizResultWithServiceNotOnline())
		return
	}

	go invokeStreamClose(streamId)
	logger.Log().Infof("发送关闭流请求完成, %s", streamId)
	render.Respond(w, r, SuccessBizResult())
}

// 要播放的流不存在事件
func StreamNotFound(w http.ResponseWriter, r *http.Request) {
	request := &message.StreamNotFoundBind{}
	if err := render.Bind(r, request); err != nil {
		logger.Log().Errorf("播放流不存在请求失败, 获取请求参数失败, %s", err.Error())
		render.Render(w, r, FailureBizResultWithParamError())
		return
	}

	if err := tcp.IsOnline(config.GetRestAddress()); err != nil {
		logger.Log().Errorf("发送拉流请求失败, %s, %s", request.Stream, err.Error())
		render.Render(w, r, FailureBizResultWithServiceNotOnline())
		return
	}

	//请求拉流
	err := invokePullStream(request)
	if err != nil {
		logger.Log().Errorf("发送拉流请求失败, %s, %s", request.Stream, err.Error())
		render.Render(w, r, FailureBizResultWithMessage(OPEARTE_STREAM_NOT_EXITS, err.Error()))
		return
	}

	logger.Log().Infof("要播放的视频流[%s]不存在,发起拉流请求成功", request.Stream)
	render.Respond(w, r, SuccessBizResult())
}

// 直播流注册/注销事件
func StreamChanged(w http.ResponseWriter, r *http.Request) {
	logger.Log().Info("stream change event received")
	render.Respond(w, r, SuccessBizResult())
}

// 直播流无人观看事件
func StreamNoneReader(w http.ResponseWriter, r *http.Request) {
	logger.Log().Info("stream no reader event received")
	render.Respond(w, r, SuccessBizResult())
}

// 发送拉流请求
func invokePullStream(request *message.StreamNotFoundBind) error {
	//根据stream 查询设备信息
	streamInfo, err := service.QueryStreamByStreamId(request.Stream)
	if err != nil || streamInfo == nil {
		logger.Log().Errorf("发送拉流请求失败, %s, %s", request.Stream, err.Error())
		return errors.New("发送拉流请求失败,未找到流的配置信息")
	}

	if len(request.App) != 0 && strings.Compare(request.App, streamInfo.AppName) != 0 {
		logger.Log().Errorf("发送拉流请求失败, 流配置信息和数据库配置信息不一致%s", request.Stream)
		return errors.New("发送拉流请求失败,未找到流的配置信息")
	}

	var builder strings.Builder
	builder.WriteString(config.GetRestUrl())
	builder.WriteString("addStreamProxy?secret=")
	builder.WriteString(config.GetServerSecret())
	builder.WriteString(streamInfo.PullStreamKey())
	logger.Log().Infof("发送拉流请求, %s", builder.String())
	rsp, err := restapi.PostUrl(builder.String())
	if err != nil {
		logger.Log().Errorf("发送拉流请求失败,%s", err.Error())
		return errors.New("发送拉流请求失败,网络请求失败")
	}

	logger.Log().Infof("发起拉流请求,%s", builder.String())
	logger.Log().Infof("发送拉流请求,接收到响应报文:%s", rsp)
	return nil
}

// 发送关闭流请求
func invokeStreamClose(streamId string) {
	//校验流编号在数据库是否已存在
	streamInfo, err := service.QueryStreamByStreamId(streamId)
	if err != nil || streamInfo == nil {
		logger.Log().Errorf("发送关闭流请求失败, %s, %s", streamId, err.Error())
		return
	}

	var builder strings.Builder
	builder.WriteString(config.GetRestUrl())
	builder.WriteString("delStreamProxy?secret=")
	builder.WriteString(config.GetServerSecret())
	builder.WriteString("&key=")
	builder.WriteString(streamInfo.CloseKey())
	logger.Log().Infof("发送关闭流请求, %s", builder.String())
	rsp, err := restapi.Get(builder.String())
	if err != nil {
		logger.Log().Errorf("发送关闭流请求失败,%s", err.Error())
	}

	logger.Log().Infof("发送关闭流请求,rsp=%s", rsp)
}

// 检查流是否在线
func invokeStreamIsOnline(streamId string) (bool, error) {
	//校验流编号在数据库是否已存在
	streamInfo, err := service.QueryStreamByStreamId(streamId)
	if err != nil || streamInfo == nil {
		logger.Log().Errorf("检查流是否在线失败, %s, %s", streamId, err.Error())
		return false, errors.New("查询流关联的设备信息失败")
	}

	var builder strings.Builder
	builder.WriteString(config.GetRestUrl())
	builder.WriteString("isMediaOnline?secret=")
	builder.WriteString(config.GetServerSecret())
	builder.WriteString("&schema=rtsp&")
	builder.WriteString(streamInfo.IsOnlineKey())
	logger.Log().Infof("发送查询流是否在线请求, %s", builder.String())
	rsp, err := restapi.Get(builder.String())
	if err != nil {
		logger.Log().Errorf("发送查询流是否在线请求失败,%s", err.Error())
	}

	request := &message.StreamIsOnlineResponse{}
	if err := json.Unmarshal([]byte(rsp), request); err != nil {
		logger.Log().Errorf("发送查询流是否在线请求失败, 获取请求参数失败, %s", err.Error())
		return false, errors.New("获取流媒体服务响应失败,请稍后重试")
	}

	logger.Log().Infof("发送查询流是否在线请求,接收到响应报文:%s", rsp)
	return request.Online, nil
}
