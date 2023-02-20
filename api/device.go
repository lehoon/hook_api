package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lehoon/hook_api/v2/library/logger"
	"github.com/lehoon/hook_api/v2/library/utils"
	"github.com/lehoon/hook_api/v2/message"
	"github.com/lehoon/hook_api/v2/service"
)

func QueryDeviceInfo(w http.ResponseWriter, r *http.Request) {
	streamId := chi.URLParam(r, "id")
	if streamId == "" {
		render.Respond(w, r, FailureBizResultWithParamError())
		return
	}

	deviceInfo, err := service.QueryDeviceByStreamId(streamId)

	if err != nil {
		logger.Log().Errorf("查询设备失败, 请稍后重试")
		render.Respond(w, r, FailureBizResultWithMessage(OPEARTE_DEVICE_NOT_EXITS, err.Error()))
		return
	}

	deviceInfo.Password = "******"
	render.Respond(w, r, SuccessBizResultWithData(deviceInfo))
	logger.Log().Info("查询设备成功")
}

func PublishDevice(w http.ResponseWriter, r *http.Request) {
	request := &message.DeviceInfoBind{}
	if err := render.Bind(r, request); err != nil {
		logger.Log().Errorf("新增设备请求失败, 获取请求参数失败, %s", err.Error())
		render.Render(w, r, FailureBizResultWithParamError())
		return
	}

	deviceInfo, _ := service.QueryDeviceByStreamId(request.StreamId)

	if deviceInfo != nil {
		logger.Log().Error("新增设备信息失败,数据已存在")
		result := FailureBizResultWithMessage(OPEARTE_DEVICE_PUBLISH_ERROR, "新增设备信息失败,数据已存在")
		render.Respond(w, r, result)
		return
	}

	if len(request.DeviceInfo.VHostName) == 0 {
		request.DeviceInfo.VHostName = "__defaultVhost__"
	}

	if !request.DeviceInfo.CanPublish() {
		logger.Log().Errorf("新增设备信息失败,数据不完整,%v", request.DeviceInfo)
		result := FailureBizResultWithMessage(OPEARTE_DEVICE_PUBLISH_ERROR, "新增设备信息失败,设备信息不完整")
		render.Respond(w, r, result)
		return
	}

	//根据streamid检查是否存在 存在返回失败信息
	if service.InsertDeviceInfo(request.DeviceInfo) != nil {
		logger.Log().Error("新增设备信息失败,请稍后重试")
		result := FailureBizResultWithMessage(OPEARTE_DEVICE_PUBLISH_ERROR, "新增设备信息失败,请稍后重试")
		render.Respond(w, r, result)
		return
	}

	logger.Log().Info("新增设备信息")
	render.Respond(w, r, SuccessBizResult())
}

func DeleteDevice(w http.ResponseWriter, r *http.Request) {
	streamId := chi.URLParam(r, "id")
	if streamId == "" {
		render.Respond(w, r, FailureBizResultWithParamError())
		return
	}

	if service.DeleteDeviceInfo(streamId) != nil {
		logger.Log().Errorf("删除设备失败, 请稍后重试")
		render.Respond(w, r, FailureBizResultWithMessage(OPEARTE_DEVICE_NOT_EXITS, "没有找到删除设备的编号,请稍后重试"))
		return
	}

	logger.Log().Info("删除设备成功")
	render.Respond(w, r, SuccessBizResult())
}

func UpdateDevice(w http.ResponseWriter, r *http.Request) {
	request := &message.DeviceInfoBind{}
	if err := render.Bind(r, request); err != nil {
		logger.Log().Errorf("更新设备请求失败, 获取请求参数失败, %s", err.Error())
		render.Render(w, r, FailureBizResultWithParamError())
		return
	}

	if request.DeviceInfo.IsEmpty() {
		logger.Log().Errorf("更新设备请求失败, 获取请求参数失败")
		render.Render(w, r, FailureBizResultWithParamError())
		return
	}

	//根据streamid检查是否存在 存在返回失败信息
	if service.InsertDeviceInfo(request.DeviceInfo) != nil {
		logger.Log().Error("更新设备请求失败,请稍后重试,%s", utils.JsonString(request.DeviceInfo))
		result := FailureBizResultWithMessage(OPEARTE_DEVICE_UPDATE_ERROR, "更新设备请求失败,请稍后重试")
		render.Respond(w, r, result)
		return
	}

	logger.Log().Info("更新设备请求成功,%s", utils.JsonString(request.DeviceInfo))
	render.Respond(w, r, SuccessBizResult())
}

// 查询设备列表
func DeviceList(w http.ResponseWriter, r *http.Request) {
	resultList, err := service.QueryDeviceList()
	if err != nil {
		logger.Log().Errorf("查询设备列表失败, %s", err.Error())
		render.Render(w, r, FailureBizResultWithDatabaseError())
		return
	}

	logger.Log().Info("查询设备列表成功")
	render.Respond(w, r, SuccessBizResultWithData(resultList))
}
