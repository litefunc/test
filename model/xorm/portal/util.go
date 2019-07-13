package model

import "encoding/json"

func Json(i interface{}) string {
	by, _ := json.Marshal(i)
	return string(by)
}
