package ext

import (
	"encoding/json"
	"net/http"
)

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
