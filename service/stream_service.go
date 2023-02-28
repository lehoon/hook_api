package service

import (
	"errors"

	"github.com/lehoon/hook_api/v2/library/database"
	"github.com/lehoon/hook_api/v2/library/logger"
	"github.com/lehoon/hook_api/v2/library/utils"
	"github.com/lehoon/hook_api/v2/message"
)

const (
	STREAM_QUERY_BY_STREAM_ID_SELECT = "select deviceid, username,password,protocol,hostname,port,vhostname,appname from device_info where streamid = ?"

	// stream operate history
	STREAM_OPERATE_HISTORY_INSERT_SQL = "insert into stream_operate_history(log_id,device_id,stream_id,server_id,client_id,opcode,message,created) values(?,?,?,?,?,?,?,datetime('now','localtime'))"
)

// 根据流编号查询流信息
func QueryStreamByStreamId(streamId string) (*message.StreamInfo, error) {
	if !database.IsOpen() {
		return nil, errors.New("数据库未打开,查询失败")
	}

	stmt, err := database.Instance().Prepare(STREAM_QUERY_BY_STREAM_ID_SELECT)
	if err != nil {
		logger.Log().Errorf("查询流信息出错,%s", err.Error())
		return nil, errors.New("查询流信息出错,请稍后重试")
	}

	defer stmt.Close()
	row := stmt.QueryRow(streamId)

	var deviceId string
	var username string
	var password string
	var protocol uint8
	var hostname string
	var port uint16
	var appname string
	var vhostname string
	err = row.Scan(&deviceId, &username, &password, &protocol, &hostname, &port, &vhostname, &appname)

	if err != nil {
		logger.Log().Errorf("查询流信息出错,%s,%s", streamId, err.Error())
		return nil, errors.New("要查询的流不存在")
	}

	return &message.StreamInfo{
		DeviceId:  deviceId,
		StreamId:  streamId,
		Username:  username,
		Password:  password,
		Protocol:  protocol,
		Hostname:  hostname,
		Port:      port,
		VHostName: vhostname,
		AppName:   appname,
	}, nil
}

// 查询所有流信息
func QueryStreamList() ([]message.StreamInfo, error) {
	resultList := []message.StreamInfo{}
	deviceList, err := QueryDeviceList()

	if err != nil {
		return resultList, nil
	}

	for _, device := range deviceList {
		resultList = append(resultList,
			message.StreamInfo{
				DeviceId:  device.DeviceId,
				StreamId:  device.StreamId,
				Username:  device.Username,
				Password:  device.Password,
				Protocol:  device.Protocol,
				Hostname:  device.Hostname,
				Port:      device.Port,
				AppName:   device.AppName,
				VHostName: device.VHostName,
			},
		)
	}

	return resultList, nil
}

// 维护流的操作信息
func StreamOperateHistory(record *StreamOperateModel) error {
	insertStmt, err := database.Instance().Prepare(STREAM_OPERATE_HISTORY_INSERT_SQL)
	if err != nil {
		logger.Log().Errorf("更新流操作数据失败,%s,%s", utils.JsonString(record), err.Error())
		return errors.New("更新流操作数据失败")
	}

	defer insertStmt.Close()

	//生成最新的log_id
	log_id := next_sequence_by_prefix_32("stream_operate_log_id")

	_, err = insertStmt.Exec(log_id, record.DeviceId, record.StreamId, record.ServerId, record.ClientId, record.Opcode, record.Message)
	if err != nil {
		logger.Log().Errorf("更新流操作数据失败,%s,%s", utils.JsonString(record), err.Error())
		return errors.New("更新流操作数据失败")
	}

	return nil
}
