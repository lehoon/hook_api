package os

import "os"

func IsPathSeparator(c uint8) bool {
	// NOTE: Windows accept / as path separator.
	return os.IsPathSeparator(c)
}

//递归创建目录
func Mkdir(path string) error {
	return os.MkdirAll(path, 0766)
}

//删除目录
func RmDir(path string) error {
	return os.RemoveAll(path)
}

//删除文件
func RemoveFile(path string) error {
	return os.RemoveAll(path)
}

//检查文件是否存在
func IsFileExist(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		return false
	}

	return true
}

//重命名文件
func RenameFile(filename, fileRename string) error {
	return os.Rename(filename, fileRename)
}

//打开文件
func OpenFile(filename string) (*os.File, error) {
	return os.Open(filename)
}

//创建新文件
func CreateFile(filename string) (*os.File, error) {
	return os.Create(filename)
}

//获取当前目录
func CurrentDir() (string, error) {
	return os.Getwd()
}

//转到目标目录
func ChangeDir(chdir string) error {
	return os.Chdir(chdir)
}
