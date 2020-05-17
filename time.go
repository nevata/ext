package ext

import (
	"time"
)

const chinaLayout = "2006-01-02 15:04:05"

//NowToStr 返回当前时间的格式化字符串
func NowToStr() string {
	return TimeToStr(time.Now())
}

//TimeToStr 返回时间的格式化字符串
func TimeToStr(t time.Time) string {
	return t.Format(chinaLayout)
}

//StrToTime 解析时间字符串
func StrToTime(s string) (time.Time, error) {
	return time.Parse(chinaLayout, s)
}
