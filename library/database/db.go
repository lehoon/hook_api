package database

import (
	"database/sql"

	"github.com/lehoon/hook_api/v2/library/logger"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dbOpenState         = false
	dbInstance  *sql.DB = nil
)

var (
	NoResultError = sql.ErrNoRows
)

func init() {
	var err error = nil
	dbInstance, err = sql.Open("sqlite3", "hook_api.db")
	if err != nil {
		logger.Log().Errorf("打开数据库文件失败, %v", err.Error())
	}

	dbOpenState = true
	logger.Log().Info("打开数据库hook_api.db成功")
}

func IsOpen() bool {
	return dbOpenState
}

func Instance() *sql.DB {
	return dbInstance
}

func Shutdown() {
	if !dbOpenState {
		return
	}

	dbInstance.Close()
}
