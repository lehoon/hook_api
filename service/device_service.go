package service

import (
	"errors"

	"github.com/lehoon/hook_api/v2/library/database"
	"github.com/lehoon/hook_api/v2/library/logger"
	"github.com/lehoon/hook_api/v2/library/utils"
	"github.com/lehoon/hook_api/v2/message"
)

func init() {
	createDeviceTable()
	logger.Log().Info("device service initialized")
}

// 根据streamid查询deviceinfo信息
func QueryDeviceByStreamId(streamId string) (*message.DeviceInfo, error) {
	if !database.IsOpen() {
		return nil, errors.New("数据库未打开,查询失败")
	}

	stmt, err := database.Instance().Prepare("select streamid,username,password,hostname,vhostname,appname from device_info where streamid = ?")
	if err != nil {
		logger.Log().Errorf("查询设备出错,%s", err.Error())
		return nil, errors.New("查询设备出错,请稍后重试")
	}

	defer stmt.Close()
	row := stmt.QueryRow(streamId)

	var streamid string
	var username string
	var password string
	var hostname string
	var appname string
	var vhostname string
	err = row.Scan(&streamid, &username, &password, &hostname, &vhostname, &appname)

	if err != nil {
		logger.Log().Errorf("查询设备数据失败,%s,%s", streamId, err.Error())
		return nil, errors.New("要查询的设备不存在")
	}

	return &message.DeviceInfo{
		StreamId:  streamid,
		Username:  username,
		Password:  password,
		Hostname:  hostname,
		VHostName: vhostname,
		AppName:   appname,
	}, nil
}

// 新增device info
func InsertDeviceInfo(deviceInfo *message.DeviceInfo) error {
	if !database.IsOpen() {
		return errors.New("数据库未打开,新增失败")
	}
	//添加新的数据
	insertSql := `insert into device_info(streamid,username,password,hostname,vhostname,appname,created) values(?,?,?,?,?,?,datetime('now','localtime'))`
	insertStmt, err := database.Instance().Prepare(insertSql)
	if err != nil {
		logger.Log().Errorf("新增设备失败,%s,%s", utils.JsonString(deviceInfo), err.Error())
		return errors.New("新增设备失败,请稍后重试")
	}

	defer insertStmt.Close()
	_, err = insertStmt.Exec(deviceInfo.StreamId, deviceInfo.Username, deviceInfo.Password, deviceInfo.Hostname, deviceInfo.VHostName, deviceInfo.AppName)
	if err != nil {
		logger.Log().Errorf("新增设备数据失败,%s,%s", utils.JsonString(deviceInfo), err.Error())
		return errors.New("新增设备失败,请稍后重试")
	}
	return nil
}

// 删除device info
func DeleteDeviceInfo(streamId string) error {
	if !database.IsOpen() {
		return errors.New("数据库未打开,删除失败")
	}
	//删除旧的数据
	deleteSql := `delete from device_info where streamid = ?`
	deleteStmt, err := database.Instance().Prepare(deleteSql)
	if err != nil {
		logger.Log().Errorf("删除设备失败,%s,%s", streamId, err.Error())
		return errors.New("删除设备失败,请稍后重试")
	}
	defer deleteStmt.Close()
	_, err = deleteStmt.Exec(streamId)
	if err != nil {
		logger.Log().Errorf("删除设备数据失败,%s, %s", streamId, err.Error())
		return errors.New("删除设备失败,请稍后重试")
	}
	return nil
}

// 更新设备信息
func UpdateDeviceInfo(deviceInfo *message.DeviceInfo) error {
	device, err := QueryDeviceByStreamId(deviceInfo.StreamId)

	if err != nil || device == nil {
		logger.Log().Error("更新设备失败, %s, %s", utils.JsonString(deviceInfo), err.Error())
		return errors.New("待更新设备数据不存在,更新设备信息失败")
	}

	if len(deviceInfo.Username) == 0 {
		deviceInfo.Username = device.Username
	}

	if len(deviceInfo.Password) == 0 {
		deviceInfo.Password = device.Password
	}

	if len(deviceInfo.Hostname) == 0 {
		deviceInfo.Hostname = device.Hostname
	}

	if len(deviceInfo.VHostName) == 0 {
		deviceInfo.VHostName = device.VHostName
	}

	if len(deviceInfo.AppName) == 0 {
		deviceInfo.AppName = device.AppName
	}

	//更新设备信息
	updateSql := "update device_info set username=?,password=?,hostname=?,vhostname=?,appname=? where streamid=?"
	updateStmt, err := database.Instance().Prepare(updateSql)
	if err != nil {
		logger.Log().Error("更新设备失败, %s, %s", utils.JsonString(deviceInfo), err.Error())
		return errors.New("更新设备信息失败,请稍后重试")
	}

	defer updateStmt.Close()
	_, err = updateStmt.Exec(deviceInfo.Username, deviceInfo.Password,
		deviceInfo.Hostname, deviceInfo.StreamId, deviceInfo.VHostName, deviceInfo.AppName)

	if err != nil {
		logger.Log().Error("更新设备失败, %s, %s", utils.JsonString(deviceInfo), err.Error())
		return errors.New("更新设备信息失败,请稍后重试")
	}

	return nil
}

// 查询所有设备信息
func QueryDeviceList() ([]message.DeviceInfo, error) {
	if !database.IsOpen() {
		return nil, errors.New("查询设备列表失败,当前数据库未建立连接")
	}

	stmt, err := database.Instance().Prepare("select streamid,username,hostname,vhostname,appname from device_info order by created desc")
	if err != nil {
		logger.Log().Error("查询设备列表失败, %s", err.Error())
		return nil, errors.New("查询设备列表失败,请稍后重试")
	}

	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		logger.Log().Error("查询设备列表失败, %s", err.Error())
		return nil, errors.New("查询设备列表失败,请稍后重试")
	}

	var streamid string
	var username string
	var hostname string
	var vhostname string
	var appname string
	resultList := []message.DeviceInfo{}

	for rows.Next() {
		err = rows.Scan(&streamid, &username, &hostname, &vhostname, &appname)
		if err == nil {
			deviceInfo := message.DeviceInfo{
				StreamId:  streamid,
				Username:  username,
				Password:  "******",
				Hostname:  hostname,
				VHostName: vhostname,
				AppName:   appname,
			}
			resultList = append(resultList, deviceInfo)
		} else {
			logger.Log().Errorf("查询设备列表出错, %s", err.Error())
		}
	}

	return resultList, nil
}

// 创建device表
func createDeviceTable() {
	if !database.IsOpen() {
		logger.Log().Info("创建设备数据库失败,当前数据库未建立连接")
		return
	}

	deviceInfoTableSql := `create table if not exists device_info (
         streamid  varchar(200) NOT NULL PRIMARY KEY,
         username  varchar(32) NOT NULL,
		 password  varchar(64) NOT NULL,
		 hostname  varchar(200) NOT NULL,
		 vhostname varchar(32) NOT NULL,
		 appname   varchar(32) NOT NULL,
         created DATE NOT NULL
         );
    `

	_, err := database.Instance().Exec(deviceInfoTableSql)
	if err != nil {
		logger.Log().Errorf("创建表结构device_info失败 %v", err)
		return
	}
}
