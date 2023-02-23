package service

import (
	"errors"

	"github.com/lehoon/hook_api/v2/library/database"
	"github.com/lehoon/hook_api/v2/library/logger"
	"github.com/lehoon/hook_api/v2/message"
)

func QueryStreamByStreamId(streamId string) (*message.StreamInfo, error) {
	if !database.IsOpen() {
		return nil, errors.New("数据库未打开,查询失败")
	}

	stmt, err := database.Instance().Prepare("select username,password,hostname,vhostname,appname from device_info where streamid = ?")
	if err != nil {
		logger.Log().Errorf("查询流信息出错,%s", err.Error())
		return nil, errors.New("查询流信息出错,请稍后重试")
	}

	defer stmt.Close()
	row := stmt.QueryRow(streamId)

	var username string
	var password string
	var hostname string
	var appname string
	var vhostname string
	err = row.Scan(&username, &password, &hostname, &vhostname, &appname)

	if err != nil {
		logger.Log().Errorf("查询流信息出错,%s,%s", streamId, err.Error())
		return nil, errors.New("要查询的流不存在")
	}

	return &message.StreamInfo{
		StreamId:  streamId,
		Username:  username,
		Password:  password,
		Hostname:  hostname,
		VHostName: vhostname,
		AppName:   appname,
	}, nil
}

func QueryStreamList() ([]message.StreamInfo, error) {
	resultList := []message.StreamInfo{}
	deviceList, err := QueryDeviceList()

	if err != nil {
		return resultList, nil
	}

	for _, device := range deviceList {
		resultList = append(resultList,
			message.StreamInfo{
				StreamId:  device.StreamId,
				Username:  device.Username,
				Password:  device.Password,
				Hostname:  device.Hostname,
				AppName:   device.AppName,
				VHostName: device.VHostName,
			},
		)
	}

	return resultList, nil
}
