package ext

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"net/url"
	"sort"
)

// Signer 签名器
type Signer struct {
	msg    url.Values
	secret string
}

// NewSigner 消息签名
func NewSigner(msg url.Values, secret string) *Signer {
	return &Signer{msg, secret}
}

// Sign 签名
func (signer *Signer) Sign() (string, error) {
	preStr := createLinkString(signer.msg)
	keyBytes := []byte(signer.secret)
	return getHmacCode(preStr, fmt.Sprintf("%x", md5.Sum(keyBytes)))
}

func createLinkString(params url.Values) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	prestr := ""
	for _, v := range keys {
		key := v
		//only use 1st value
		value := params[key][0]
		if prestr != "" {
			prestr += "&"
		}
		prestr += fmt.Sprintf("%s=%s", key, value)
	}
	return prestr
}

func getHmacCode(value, key string) (string, error) {
	mac := hmac.New(sha256.New, []byte(key))
	if _, err := io.WriteString(mac, value); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", mac.Sum(nil)), nil
}
