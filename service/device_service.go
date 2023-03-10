package service

import (
	"errors"
	"strings"

	"github.com/lehoon/hook_api/v2/library/database"
	"github.com/lehoon/hook_api/v2/library/logger"
	"github.com/lehoon/hook_api/v2/library/utils"
	"github.com/lehoon/hook_api/v2/message"
)

func init() {
	logger.Log().Info("device service initialized")
}

// 根据deviceid查询deviceinfo信息
func QueryDeviceByDeviceId(deviceid string) (*message.DeviceInfo, error) {
	if !database.IsOpen() {
		return nil, errors.New("数据库未打开,查询失败")
	}

	stmt, err := database.Instance().Prepare("select streamid,username,password,protocol,hostname,port,vhostname,appname from device_info where deviceid = ?")
	if err != nil {
		logger.Log().Errorf("查询设备出错,%s", err.Error())
		return nil, errors.New("查询设备出错,请稍后重试")
	}

	defer stmt.Close()
	row := stmt.QueryRow(deviceid)

	var streamid string
	var username string
	var password string
	var protocol uint8
	var hostname string
	var port uint16
	var appname string
	var vhostname string
	err = row.Scan(&streamid, &username, &password, &protocol, &hostname, &port, &vhostname, &appname)

	if err != nil {
		logger.Log().Errorf("查询设备数据失败,%s,%s", deviceid, err.Error())
		return nil, errors.New("要查询的设备不存在")
	}

	return &message.DeviceInfo{
		DeviceId:  deviceid,
		StreamId:  streamid,
		Username:  username,
		Password:  password,
		Protocol:  protocol,
		Hostname:  hostname,
		Port:      port,
		VHostName: vhostname,
		AppName:   appname,
	}, nil
}

// 查询设备表数量
func count_device_table() int {
	if !database.IsOpen() {
		return 0
	}

	stmt, err := database.Instance().Prepare("select count(*) as totalcount from device_info")
	if err != nil {
		logger.Log().Errorf("检查设备表失败,%s", err.Error())
		return 0
	}

	defer stmt.Close()
	row := stmt.QueryRow()

	var totalcount int
	err = row.Scan(&totalcount)

	if err != nil {
		return 0
	}

	return totalcount
}

// 查询流最大编号
func device_streamid_max() (int32, error) {
	if !database.IsOpen() {
		return 0, errors.New("查询设备表流序号失败,当前数据库未建立连接")
	}

	if count_device_table() == 0 {
		return 0, nil
	}

	stmt, err := database.Instance().Prepare("select max(cast(streamid as int)) as streamid from device_info")
	if err != nil {
		logger.Log().Error("查询设备表流序号失败, %s", err.Error())
		return 0, errors.New("查询设备表流序号失败,请稍后重试")
	}

	defer stmt.Close()
	row := stmt.QueryRow()

	var streamid string
	err = row.Scan(&streamid)

	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		logger.Log().Errorf("生成流序列号失败,未初始化数据,%s", err.Error())
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return string_to_int32(streamid), nil
}

// 新增device info
func InsertDeviceInfo(deviceInfo *message.DeviceInfo) error {
	if !database.IsOpen() {
		return errors.New("数据库未打开,新增失败")
	}

	streamid_max, err := device_streamid_max()
	if err != nil {
		logger.Log().Errorf("新增设备失败,数据库发送错误, %s,%s", utils.JsonString(deviceInfo), err.Error())
		return errors.New("新增设备失败,请稍后重试")
	}

	streamid, err := next_sequece()
	if err != nil {
		logger.Log().Errorf("新增设备失败,获取流序列号失败, %s,%s", utils.JsonString(deviceInfo), err.Error())
		return errors.New("新增设备失败,请稍后重试")
	}

	streamid_new_int := string_to_int32(streamid)

	if streamid_new_int <= streamid_max {
		streamid_new_int = streamid_max + 1
		streamid, err = update_sequence(streamid_new_int)
		if err != nil {
			logger.Log().Errorf("新增设备失败,获取流序列号失败, %s,%s", utils.JsonString(deviceInfo), err.Error())
			return errors.New("新增设备失败,请稍后重试")
		}
	}

	if deviceInfo.Protocol == 0 {
		deviceInfo.Protocol = message.STREAM_TRANSPORT_PROTOCOL_RTSP
	}

	deviceInfo.Port = message.GetProtocolPort(deviceInfo.Protocol, deviceInfo.Port)

	//添加新的数据
	insertSql := `insert into device_info(deviceid,streamid,username,password,protocol,hostname,port,vhostname,appname,created) values(?,?,?,?,?,?,?,?,?,datetime('now','localtime'))`
	insertStmt, err := database.Instance().Prepare(insertSql)
	if err != nil {
		logger.Log().Errorf("新增设备失败,%s,%s", utils.JsonString(deviceInfo), err.Error())
		return errors.New("新增设备失败,请稍后重试")
	}

	defer insertStmt.Close()
	_, err = insertStmt.Exec(deviceInfo.DeviceId, streamid, deviceInfo.Username, deviceInfo.Password,
		deviceInfo.Protocol, deviceInfo.Hostname, deviceInfo.Port, deviceInfo.VHostName, deviceInfo.AppName)
	if err != nil {
		logger.Log().Errorf("新增设备数据失败,%s,%s", utils.JsonString(deviceInfo), err.Error())
		return errors.New("新增设备失败,请稍后重试")
	}
	return nil
}

// 删除device info
func DeleteDeviceInfo(deviceId string) error {
	if !database.IsOpen() {
		return errors.New("数据库未打开,删除失败")
	}
	//删除旧的数据
	deleteSql := `delete from device_info where streamid = ?`
	deleteStmt, err := database.Instance().Prepare(deleteSql)
	if err != nil {
		logger.Log().Errorf("删除设备失败,%s,%s", deviceId, err.Error())
		return errors.New("删除设备失败,请稍后重试")
	}
	defer deleteStmt.Close()
	_, err = deleteStmt.Exec(deviceId)
	if err != nil {
		logger.Log().Errorf("删除设备数据失败,%s, %s", deviceId, err.Error())
		return errors.New("删除设备失败,请稍后重试")
	}
	return nil
}

// 更新设备信息
func UpdateDeviceInfo(deviceInfo *message.DeviceInfo) error {
	device, err := QueryDeviceByDeviceId(deviceInfo.DeviceId)

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
	updateSql := "update device_info set username=?,password=?,hostname=?,vhostname=?,appname=? where deviceid=?"
	updateStmt, err := database.Instance().Prepare(updateSql)
	if err != nil {
		logger.Log().Error("更新设备失败, %s, %s", utils.JsonString(deviceInfo), err.Error())
		return errors.New("更新设备信息失败,请稍后重试")
	}

	defer updateStmt.Close()
	_, err = updateStmt.Exec(deviceInfo.Username, deviceInfo.Password,
		deviceInfo.Hostname, deviceInfo.DeviceId, deviceInfo.VHostName, deviceInfo.AppName)

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

	stmt, err := database.Instance().Prepare("select deviceid,username,protocol,hostname,port,vhostname,appname from device_info order by created desc")
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

	var deviceid string
	var username string
	var protocol uint8
	var hostname string
	var port uint16
	var vhostname string
	var appname string
	resultList := []message.DeviceInfo{}

	for rows.Next() {
		err = rows.Scan(&deviceid, &username, &protocol, &hostname, &port, &vhostname, &appname)
		if err == nil {
			deviceInfo := message.DeviceInfo{
				DeviceId:  deviceid,
				Username:  username,
				Password:  "******",
				Protocol:  protocol,
				Hostname:  hostname,
				Port:      port,
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
