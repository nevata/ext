package ext

import (
	"os"
	"os/exec"
	"path/filepath"
)

//FileExist 文件和目录判断
func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

//GetAppPath 获取应用程序路径
func GetAppPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	return filepath.Dir(path), nil
}

//MustGetAppPath 获取应用程序路径
func MustGetAppPath() string {
	s, e := GetAppPath()
	if e != nil {
		panic(e)
	}
	return s
}

//JSONTimeToStr 转换jsontime对象指针，空指针返回空字符串
func JSONTimeToStr(jsonTime *JSONTime) string {
	if jsonTime != nil {
		return TimeToStr(jsonTime.Time)
	} else {
		return ""
	}
}
