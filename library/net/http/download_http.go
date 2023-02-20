package http

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

//http下载文件
func DonwloadFile(url, path string) (string, error) {
	fileName := filepath.Base(url)
	newFile  := path + fileName

	res, err := http.Get(url)

	if err != nil {
		return newFile, err
	}

	file, err := os.Create(newFile)

	if err != nil {
		return newFile, err
	}

	io.Copy(file, res.Body)
	return newFile, nil
}
