package ext

import "time"

//NowToStr 返回格式化的时间字符串
func NowToStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
