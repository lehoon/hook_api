package service

import (
	"errors"

	"github.com/lehoon/hook_api/v2/library/database"
	"github.com/lehoon/hook_api/v2/library/logger"
)

const (
	STREAM_HEARTBEAT_INSET_OR_UPDATE_SQL = "insert or replace into stream_server_heartbeat (server_id,created) values(?, datetime('now','localtime'))"
	STREAM_SERVER_OPERATE_INSET_SQL      = "insert into stream_server_operate (log_id,server_id,opcode,messagecreated) values(?,?,?,?,datetime('now','localtime'))"
)

// 更新服务器操作数据
func StreamServerOperate(serverId, opcode, message string) error {
	insertStmt, err := database.Instance().Prepare(STREAM_SERVER_OPERATE_INSET_SQL)
	if err != nil {
		logger.Log().Errorf("更新服务器操作数据失败,%s,%s", serverId, err.Error())
		return errors.New("更新服务器操作数据失败")
	}

	defer insertStmt.Close()

	//生成最新的log_id
	log_id := next_sequence_by_prefix_32("stream_server_log_id")

	_, err = insertStmt.Exec(log_id, serverId, opcode, message)
	if err != nil {
		logger.Log().Errorf("更新服务器操作数据失败,%s,%s", serverId, err.Error())
		return errors.New("更新服务器操作数据失败")
	}

	return nil
}

// 更新服务器心跳数据
func StreamServerHeartBeat(serverId string) error {
	insertStmt, err := database.Instance().Prepare(STREAM_HEARTBEAT_INSET_OR_UPDATE_SQL)
	if err != nil {
		logger.Log().Errorf("更新服务器心跳数据失败,%s,%s", serverId, err.Error())
		return errors.New("更新服务器心跳数据失败")
	}

	defer insertStmt.Close()
	_, err = insertStmt.Exec(serverId)
	if err != nil {
		logger.Log().Errorf("更新服务器心跳数据失败,%s,%s", serverId, err.Error())
		return errors.New("更新服务器心跳数据失败")
	}

	return nil
}
