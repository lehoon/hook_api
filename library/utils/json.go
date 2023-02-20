package utils

import "encoding/json"

func JsonString(obj interface{}) string {
	buf, err := json.Marshal(obj)
	if err != nil {
		return ""
	}

	return string(buf)
}
