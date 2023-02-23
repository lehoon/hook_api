package service

import (
	"errors"
	"strconv"
	"strings"

	"github.com/lehoon/hook_api/v2/library/database"
	"github.com/lehoon/hook_api/v2/library/logger"
)

const DEFAULT_STREAM_SEQUENCE_ID = 10001

// 获取流序号
func next_sequece() (string, error) {
	if !database.IsOpen() {
		return "", errors.New("数据库未打开,生成流序列号失败")
	}

	var sequence_no int32
	seqno, err := current_sequence()

	if err != nil {
		return "", err
	}

	// TODO:生成新的流序列号 并更新数据库
	sequence_no1 := string_to_int32(seqno)

	if sequence_no1 == 0 {
		sequence_no = DEFAULT_STREAM_SEQUENCE_ID
	} else {
		sequence_no = sequence_no1
	}

	if sequence_no == DEFAULT_STREAM_SEQUENCE_ID {
		create_sequence(sequence_no)
	}

	//写入数据库
	return update_sequence(sequence_no)
}

func create_sequence(sequence_no int32) error {
	updateSql := "insert into stream_sequence(sequenceid) values(?)"
	updateStmt, err := database.Instance().Prepare(updateSql)
	if err != nil {
		logger.Log().Error("创建流序号表失败, %d, %s", sequence_no, err.Error())
		return errors.New("创建流序号表失败,请稍后重试")
	}

	defer updateStmt.Close()
	_, err = updateStmt.Exec(sequence_no)

	if err != nil {
		logger.Log().Error("创建流序号表失败, %d, %s", sequence_no, err.Error())
		return errors.New("创建流序号表失败,请稍后重试")
	}

	return nil
}

func update_sequence(sequence_no int32) (string, error) {
	//更新流序号
	sequence_no_new := sequence_no + 1
	updateSql := "update stream_sequence set sequenceid=? where sequenceid=?"
	updateStmt, err := database.Instance().Prepare(updateSql)
	if err != nil {
		logger.Log().Error("更新流序号表失败, %d, %s", sequence_no, err.Error())
		return "", errors.New("更新流序号表失败,请稍后重试")
	}

	defer updateStmt.Close()
	_, err = updateStmt.Exec(sequence_no_new, sequence_no)

	if err != nil {
		logger.Log().Error("更新流序号表失败, %d, %s", sequence_no, err.Error())
		return "", errors.New("更新流序号表失败,请稍后重试")
	}

	logger.Log().Error("更新流序号表成功, %d", sequence_no_new)
	return int32_to_string(sequence_no_new), nil
}

// 应急更新序列号  防止因为序列号错乱导致写入多条数据
func update_sequence_v1(sequence_no, sequence_no_old int32) (string, error) {
	//更新流序号
	updateSql := "update stream_sequence set sequenceid=? where sequenceid=?"
	updateStmt, err := database.Instance().Prepare(updateSql)
	if err != nil {
		logger.Log().Error("更新流序号表失败, %d, %d, %s", sequence_no, sequence_no_old, err.Error())
		return "", errors.New("更新流序号表失败,请稍后重试")
	}

	defer updateStmt.Close()
	_, err = updateStmt.Exec(sequence_no, sequence_no_old)

	if err != nil {
		logger.Log().Error("更新流序号表失败, %d, %d, %s", sequence_no, sequence_no_old, err.Error())
		return "", errors.New("更新流序号表失败,请稍后重试")
	}

	logger.Log().Error("更新流序号表成功, %d, %d", sequence_no, sequence_no_old)
	return int32_to_string(sequence_no), nil
}

func current_sequence() (string, error) {
	if !database.IsOpen() {
		return "", errors.New("数据库未打开,生成流序列号失败")
	}

	stmt, err := database.Instance().Prepare("select sequenceid from stream_sequence")
	if err != nil {
		logger.Log().Errorf("生成流序列号失败,%s", err.Error())
		return "", errors.New("生成流序列号失败,请稍后重试")
	}

	defer stmt.Close()
	row := stmt.QueryRow()

	var sequenceid string
	err = row.Scan(&sequenceid)

	logger.Log().Infof("查询到的流序列为%s", sequenceid)

	strings.TrimSpace(sequenceid)

	if len(sequenceid) > 0 {
		return sequenceid, nil
	}

	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		logger.Log().Errorf("生成流序列号失败,未初始化数据,%s", err.Error())
		return strconv.FormatInt(int64(DEFAULT_STREAM_SEQUENCE_ID), 10), nil
	}

	logger.Log().Errorf("生成流序列号失败,%s", err.Error())
	return "", err
}

// reapir stream sequence table
func repair_stream_sequence_table() {
	if !database.IsOpen() {
		return
	}

	stmt, err := database.Instance().Prepare("select count(*) as totalcount from stream_sequence")
	if err != nil {
		logger.Log().Errorf("检查视频流序列号表失败,%s", err.Error())
		return
	}

	defer stmt.Close()
	row := stmt.QueryRow()

	var totalcount int
	err = row.Scan(&totalcount)

	if err != nil {
		return
	}

	if totalcount <= 1 {
		return
	}

	deleteSql := `delete from stream_sequence`
	deleteStmt, err := database.Instance().Prepare(deleteSql)
	if err != nil {
		logger.Log().Errorf("删除视频流序列号表,%s", err.Error())
		return
	}

	defer deleteStmt.Close()
	_, err = deleteStmt.Exec()
	if err != nil {
		logger.Log().Errorf("删除视频流序列号表,%s", err.Error())
		return
	}
}

func string_to_int32(s string) int32 {
	ivalue, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}

	return int32(ivalue)
}

func int32_to_string(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}
