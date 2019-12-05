package ext

import (
	"encoding/json"
	"net/http"
)

func HandleError(w http.ResponseWriter, e error) {
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"result":  false,
		"message": e.Error(),
	})

	if err != nil {
		panic(err)
	}
}

func HandleSuccess(w http.ResponseWriter, data interface{}) {
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"result": true,
		"data":   data,
	})

	if err != nil {
		panic(err)
	}
}
