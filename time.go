package ext

import (
	"time"
)

const timeLayout = "2006-01-02 15:04:05"
const dateLayout = "2006-01-02"

//Now 返回当前的UTC时间(精确到秒)
func Now() time.Time {
	t := time.Now().UTC()
	y, m, d := t.Date()
	h, n, s := t.Clock()
	return time.Date(y, m, d, h, n, s, 0, t.Location())
}

//NowToStr 返回当前时间的格式化字符串
func NowToStr() string {
	return TimeToStr(Now())
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
