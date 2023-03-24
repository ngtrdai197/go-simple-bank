package util

import "encoding/json"

func Stringify(data interface{}) string {
	marshal, _ := json.Marshal(data)
	return string(marshal)
}
