package ext

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
)

// FileExist 文件判断
func FileExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}

// DirExist 目录判断
func DirExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	if !info.IsDir() {
		return false
	}
	return true
}

// GetAppPath 获取应用程序路径
func GetAppPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

// MustGetAppPath 获取应用程序路径
func MustGetAppPath() string {
	s, e := GetAppPath()
	if e != nil {
		panic(e)
	}
	return s
}

// JSONTimeToStr 转换jsontime对象指针，空指针返回空字符串
func JSONTimeToStr(jsonTime *JSONTime) string {
	if jsonTime != nil {
		return TimeToStr(jsonTime.Time)
	}

	return ""
}

// VerifyEmailFormat 电子邮箱格式检测
func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// PrintErr 输出错误以及堆栈信息
func PrintErr(err error) {
	log.Println(err)
	debug.PrintStack()
}

// HandlePassword 加密
func HandlePassword(pwd, seed string) (s string) {
	s = fmt.Sprintf("%s&%s", pwd, seed)
	s = fmt.Sprintf("%x", md5.Sum([]byte(s)))
	return
}
