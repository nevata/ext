package ext

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

//SID sessionid
var SID string

//HandleError 错误返回
func HandleError(w http.ResponseWriter, e error) {
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"result":  false,
		"message": e.Error(),
	})

	if err != nil {
		panic(err)
	}
}

//HandleSuccess 成功返回
func HandleSuccess(w http.ResponseWriter, data interface{}) {
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"result": true,
		"data":   data,
	})

	if err != nil {
		panic(err)
	}
}

//Post 用json格式发送
func Post(url string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	if SID != "" {
		req.Header.Set("Authorization", fmt.Sprintf("DSSESSION %s", SID))
	}
	return http.DefaultClient.Do(req)
}
