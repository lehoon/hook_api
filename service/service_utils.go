package service

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

func current_datetime() string {
	current_time := time.Now()
	return current_time.Format("2006-01-02 15:04:05")
}

// 生成32位md5序号  通过prefix+时间 生成md5数据
func next_sequence_by_prefix_32(prefix string) string {
	current_time := current_datetime()
	prefix += current_time
	buff := md5.Sum([]byte(prefix))
	return hex.EncodeToString(buff[:])
}
