package service

import (
	"github.com/lehoon/hook_api/v2/library/database"
	"github.com/lehoon/hook_api/v2/library/logger"
)

func init() {
	createDeviceTable()
	createStreamSequenceTable()
	repair_stream_sequence_table()
	createStreamOperateHistoryTable()
	createMediaServerOperateHistoryTable()
	createMediaServerHeartBeatTable()
	logger.Log().Info("database service initialized")
}

// 创建device表
func createDeviceTable() {
	if !database.IsOpen() {
		logger.Log().Info("创建设备表失败,当前数据库未建立连接")
		panic("创建设备表失败,当前数据库未建立连接")
	}

	deviceInfoTableSql := `create table if not exists device_info (
         deviceid  varchar(32) NOT NULL PRIMARY KEY,
		 streamid  varchar(32) NOT NULL,
         username  varchar(32) NOT NULL,
		 password  varchar(64) NOT NULL,
		 protocol  integer NOT NULL,
		 hostname  varchar(200) NOT NULL,
		 port      integer NOT NULL,
		 vhostname varchar(32) NOT NULL,
		 appname   varchar(32) NOT NULL,
         created   DATE NOT NULL
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
		logger.Log().Info("创建自增序列表失败,当前数据库未建立连接")
		panic("创建自增序列表失败,当前数据库未建立连接")
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

// 创建流操作播放历史表
func createStreamOperateHistoryTable() {
	if !database.IsOpen() {
		logger.Log().Info("创建流操作历史表失败,当前数据库未建立连接")
		panic("创建流操作历史表失败,当前数据库未建立连接")
	}

	/**
	device_id  设备编号  关联device_info.device_id
	stream_id  流编号    关联device_info.stream_id   冗余设计 防止通过表关联统计数据
	server_id  该流在哪台服务器上
	client_id  客户端是什么类型    flv  mp4 rtsp/vlc  需要预定义该有效值
	opcode     操作类型   open close play pause continue
	*/
	operateHistoryTableSql := `create table if not exists stream_operate_history (
		log_id    varchar(32) NOT NULL PRIMARY KEY,
		device_id varchar(32) NOT NULL,
		stream_id varchar(32) NOT NULL,
		server_id varchar(32) NOT NULL,
		client_id varchar(32) NOT NULL,
		opcode    varchar(16) NOT NULL,
		message   varchar(255) NOT NULL,
		created   DATE NOT NULL
        );
    `

	_, err := database.Instance().Exec(operateHistoryTableSql)
	if err != nil {
		logger.Log().Errorf("创建流操作历史表stream_operate_history失败 %v", err)
		return
	}
}

// 创建多媒体服务器操作历史表
func createMediaServerOperateHistoryTable() {
	if !database.IsOpen() {
		logger.Log().Info("创建多媒体服务器操作历史表失败,当前数据库未建立连接")
		panic("创建多媒体服务器操作历史表失败,当前数据库未建立连接")
	}

	operateHistoryTableSql := `create table if not exists stream_server_operate (
		log_id    varchar(32) NOT NULL PRIMARY KEY,
		server_id varchar(32) NOT NULL,
		opcode    varchar(16) NOT NULL,
		message   varchar(255) NOT NULL,
		created   DATE NOT NULL
        );
    `

	_, err := database.Instance().Exec(operateHistoryTableSql)
	if err != nil {
		logger.Log().Errorf("创建表结构stream_server_operate失败 %v", err)
		return
	}
}

// 创建多媒体服务器心跳表
func createMediaServerHeartBeatTable() {
	if !database.IsOpen() {
		logger.Log().Info("创建多媒体服务器心跳表失败,当前数据库未建立连接")
		panic("创建多媒体服务器心跳表失败,当前数据库未建立连接")
	}

	heartbeatTableSql := `create table if not exists stream_server_heartbeat (
		server_id varchar(32) NOT NULL PRIMARY KEY,
		created   DATE NOT NULL
        );
    `

	_, err := database.Instance().Exec(heartbeatTableSql)
	if err != nil {
		logger.Log().Errorf("创建表结构stream_server_heartbeat失败 %v", err)
		return
	}
}
