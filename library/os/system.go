package os

import (
	"os"
	"runtime"
	"strings"
)

//获取系统类型  linux/window
func GetOSName() string {
	return runtime.GOOS
}

//是否有java环境变量
func HasJavaEnv() bool {
	javaEnv := os.Getenv("JAVA_HOME")
	javaEnv = strings.TrimSpace(javaEnv)
	return len(javaEnv) > 0
}

func Env() []string{
	return os.Environ()
}

func SetEnv(key, value string) {
	os.Setenv(key, value)
}
