package api

import (
	"net/http"

	"github.com/go-chi/render"
)

const (
	OPEARTE_SUCCESS = 0 //操作成功

	OPEARTE_REQUEST_PARAM_ERROR = 101 //请求参数错误

	OPEARTE_DEVICE_NOT_EXITS     = 201 //设备不存在
	OPEARTE_DEVICE_PUBLISH_ERROR = 202 //新增设备失败
	OPEARTE_DEVICE_UPDATE_ERROR  = 203 //更新设备失败

	OPEARTE_SERVICE_NOT_ONLINE = 301 //多媒体服务未启动

	OPEARTE_DATABASE_ERROR = 401 //数据库异常

	OPEARTE_STREAM_NOT_EXITS      = 501 //流不存在
	OPEARTE_STREAM_ISONLINE_ERROR = 502 //检查流是否在线失败
)

type BusinessCodeMsg struct {
	Code    uint16 `json:"code"`
	Message string `json:"message"`
}

var (
	businessCodeMsgData []BusinessCodeMsg
)

func init() {
	businessCodeMsgData = append(businessCodeMsgData, BusinessCodeMsg{
		Code:    OPEARTE_SUCCESS,
		Message: "操作成功",
	})

	businessCodeMsgData = append(businessCodeMsgData, BusinessCodeMsg{
		Code:    OPEARTE_REQUEST_PARAM_ERROR,
		Message: "请求参数错误",
	})

	businessCodeMsgData = append(businessCodeMsgData, BusinessCodeMsg{
		Code:    OPEARTE_DEVICE_NOT_EXITS,
		Message: "设备不存在",
	})

	businessCodeMsgData = append(businessCodeMsgData, BusinessCodeMsg{
		Code:    OPEARTE_DEVICE_PUBLISH_ERROR,
		Message: "新增设备失败",
	})

	businessCodeMsgData = append(businessCodeMsgData, BusinessCodeMsg{
		Code:    OPEARTE_DEVICE_UPDATE_ERROR,
		Message: "更新设备失败",
	})

	businessCodeMsgData = append(businessCodeMsgData, BusinessCodeMsg{
		Code:    OPEARTE_SERVICE_NOT_ONLINE,
		Message: "多媒体服务未启动",
	})

	businessCodeMsgData = append(businessCodeMsgData, BusinessCodeMsg{
		Code:    OPEARTE_DATABASE_ERROR,
		Message: "数据库异常",
	})

	businessCodeMsgData = append(businessCodeMsgData, BusinessCodeMsg{
		Code:    OPEARTE_STREAM_NOT_EXITS,
		Message: "流不存在",
	})

	businessCodeMsgData = append(businessCodeMsgData, BusinessCodeMsg{
		Code:    OPEARTE_STREAM_ISONLINE_ERROR,
		Message: "检查流是否在线失败",
	})
}

func OperateCodeMessage() []BusinessCodeMsg {
	return businessCodeMsgData
}

type BusinessResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (br *BusinessResult) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, br.Code)
	return nil
}

func SuccessBizResult() *BusinessResult {
	return &BusinessResult{
		Code:    OPEARTE_SUCCESS,
		Message: "操作成功",
		Data:    nil,
	}
}

func SuccessBizResultWithData(data interface{}) *BusinessResult {
	return &BusinessResult{
		Code:    OPEARTE_SUCCESS,
		Message: "操作成功",
		Data:    data,
	}
}

func FailureBizResult() *BusinessResult {
	return FailureBizResultWithMessage(100, "操作失败")
}

func FailureBizResultWithCode(code int) *BusinessResult {
	return FailureBizResultWithMessage(code, "操作失败")
}

func FailureBizResultWithMessage(code int, message string) *BusinessResult {
	return &BusinessResult{
		Code:    code,
		Message: "操作失败",
		Data:    message,
	}
}

// 参数错误异常信息
func FailureBizResultWithParamError() *BusinessResult {
	return FailureBizResultWithMessage(OPEARTE_REQUEST_PARAM_ERROR, "请求参数错误")
}

// 设备不存在
func FailureBizResultWithDeviceNotExist() *BusinessResult {
	return FailureBizResultWithMessage(OPEARTE_DEVICE_NOT_EXITS, "设备不存在")
}

// 流媒体服务未启动
func FailureBizResultWithServiceNotOnline() *BusinessResult {
	return FailureBizResultWithMessage(OPEARTE_SERVICE_NOT_ONLINE, "流媒体服务未启动")
}

// 数据库异常
func FailureBizResultWithDatabaseError() *BusinessResult {
	return FailureBizResultWithMessage(OPEARTE_DATABASE_ERROR, "数据库异常")
}

func FailureBizResultWithStreamNotExists() *BusinessResult {
	return FailureBizResultWithMessage(OPEARTE_STREAM_NOT_EXITS, "流不存在")
}
