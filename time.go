package ext

import (
	"time"
)

const timeLayout = "2006-01-02 15:04:05"
const dateLayout = "2006-01-02"

//NowToStr 返回当前时间的格式化字符串
func NowToStr() string {
	return TimeToStr(time.Now())
}

//TimeToStr 返回时间的格式化字符串
func TimeToStr(t time.Time) string {
	return t.Format(timeLayout)
}

//DateToStr 返回日期的格式化字符串
func DateToStr(d time.Time) string {
	return d.Format(dateLayout)
}

//StrToTime 解析时间字符串
func StrToTime(s string) (time.Time, error) {
	return time.Parse(timeLayout, s)
}

//StrToDate 解析日期字符串
func StrToDate(s string) (time.Time, error) {
	return time.Parse(dateLayout, s)
}
