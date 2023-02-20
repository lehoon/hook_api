package os

import (
	"os"
	"path/filepath"
)

//检查java环境  true  已经安装  false 没有安装
func IsJavaInstall() bool {
	//查询env变量是否有JAVA_HOME
	javaHome := os.Getenv("JAVA_HOME")
	if len(javaHome) == 0 {
		return false
	}

	javaBinPath := javaHome + string(filepath.Separator) + "bin" + string(filepath.Separator) + "java"
	if IsFileExist(javaBinPath) {
		return true
	}

	//根据JAVA_HOME找java.exe
	return false
}

//检查apama环境  true  已经安装  false 没有安装
func IsApamaInstall() bool {
	//查询env变量是否有APAMA_HOME
	apamaHome := os.Getenv("APAMA_HOME")
	if len(apamaHome) == 0 {
		return false
	}

	apamaBinPath := apamaHome + string(filepath.Separator) + "bin" + string(filepath.Separator) + "correlator"
	if IsFileExist(apamaBinPath) {
		return true
	}

	//根据APAMA_HOME找correlator.exe
	return false
}

