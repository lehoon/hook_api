package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

func Get(url string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	rsp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	defer rsp.Body.Close()
	result, _ := io.ReadAll(rsp.Body)
	return string(result), nil
}

func DeleteUrl(url string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", err
	}

	//req.Header.Set("Content-Type", contentType)
	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer rsp.Body.Close()
	result, _ := io.ReadAll(rsp.Body)
	return string(result), nil
}

func PostUrl(url string) (string, error) {
	var data map[string][]string
	return PostWithoutBody(url, data)
}

func PostWithoutBody(url string, data url.Values) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	rsp, err := client.PostForm(url, data)
	if err != nil {
		return "", err
	}

	defer rsp.Body.Close()
	result, _ := io.ReadAll(rsp.Body)
	return string(result), nil
}

func PostWithBody(url string, data interface{}, contentType string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	json_data, _ := json.Marshal(data)
	rsp, err := client.Post(url, contentType, bytes.NewBuffer(json_data))

	if err != nil {
		return "", err
	}

	defer rsp.Body.Close()
	result, _ := io.ReadAll(rsp.Body)
	return string(result), nil
}
