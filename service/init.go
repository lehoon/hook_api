package service

import (
	"github.com/lehoon/hook_api/v2/library/database"
	"github.com/lehoon/hook_api/v2/library/logger"
)

func init() {
	createDeviceTable()
	createStreamSequenceTable()
	repair_stream_sequence_table()
	logger.Log().Info("database service initialized")
}

// 创建device表
func createDeviceTable() {
	if !database.IsOpen() {
		logger.Log().Info("创建设备数据库失败,当前数据库未建立连接")
		panic("创建设备数据库失败,当前数据库未建立连接")
	}

	deviceInfoTableSql := `create table if not exists device_info (
         deviceid  varchar(200) NOT NULL PRIMARY KEY,
		 streamid  varchar(32) NOT NULL,
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

// 创建自增序列表维护一个流序号
func createStreamSequenceTable() {
	if !database.IsOpen() {
		logger.Log().Info("创建流序列数据库失败,当前数据库未建立连接")
		panic("创建流序列数据库失败,当前数据库未建立连接")
	}

	sequenceTableSql := `create table if not exists stream_sequence (
         sequenceid integer );
    `

	_, err := database.Instance().Exec(sequenceTableSql)
	if err != nil {
		logger.Log().Errorf("创建表结构stream_sequence失败 %v", err)
		return
	}
}
