package ext

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// SID sessionid
var SID string

// ErrInternal 内部错误
var ErrInternal = errors.New("内部错误")

// ErrSign 签名不正确
var ErrSign = errors.New("签名不正确")

// HandleExcept 处理异常
func HandleExcept(w http.ResponseWriter, e error) {
	HandleError(w, ErrInternal)
	PrintErr(e)
}

// HandleError 错误返回
func HandleError(w http.ResponseWriter, e error) {
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    -1,
		"message": e.Error(),
	})

	if err != nil {
		panic(err)
	}
}

// HandleSuccess 成功返回
func HandleSuccess(w http.ResponseWriter, data interface{}) {
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 0,
		"data": data,
	})

	if err != nil {
		panic(err)
	}
}

// HandleMessage 检查错误并返回data
func HandleMessage(resp *http.Response) (json.RawMessage, error) {
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	var msg struct {
		Result bool            `json:"result"`
		Msg    string          `json:"message"`
		Data   json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, err
	}
	if msg.Result == false {
		return nil, errors.New(msg.Msg)
	}
	return msg.Data, nil
}

// HandleSID 处理返回的sid
func HandleSID(data json.RawMessage) error {
	var sess struct {
		SID string `json:"sid"`
	}
	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&sess); err != nil {
		return err
	}
	SID = sess.SID
	return nil
}

// Post 用json格式发送
func Post(url string, o interface{}) (resp *http.Response, err error) {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(o); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	if SID != "" {
		req.Header.Set("Authorization",
			fmt.Sprintf("DSSESSION %s", base64.StdEncoding.EncodeToString([]byte(SID))))
	}
	return http.DefaultClient.Do(req)
}

// FormValue 获取请求参数
func FormValue(r *http.Request, key string) (string, bool) {
	v := r.FormValue(key)
	if v == "" {
		_, b := r.Form[key]
		return v, b
	}
	return v, true
}

// CheckSign API接口校验签名
func CheckSign(exchange func(apiKey string) string, r *http.Request) error {
	sign := r.FormValue("sign")
	r.Form.Del("sign")

	apiKey := r.FormValue("api_key")
	secretKey := exchange(apiKey)

	val, err := NewSigner(r.Form, secretKey).Sign()
	if err != nil {
		return err
	}

	if sign != val {
		return ErrSign
	}

	return nil
}
