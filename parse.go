package ext

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

//checkFilterStr 校验请求的过滤参数
func checkFilterStr(filterStr string) error {
	matched, err := regexp.MatchString("^AND \\(.*\\)", filterStr)
	if err != nil {
		return err
	}

	if !matched {
		return fmt.Errorf("filterStr参数不正确")
	}

	return nil
}

//Parse2M 读取请求参数的pageSize,pageIndex两个参数
func Parse2M(r *http.Request) (pageIndex, pageSize int, err error) {
	if err := r.ParseForm(); err != nil {
		return 0, 0, err
	}

	pageSize, err = strconv.Atoi(r.FormValue("pageSize"))
	if err != nil {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	pageIndex, err = strconv.Atoi(r.FormValue("pageIndex"))
	if err != nil {
		pageIndex = 0
	}

	return pageIndex, pageSize, nil
}

//Parse3M 读取请求参数的pageSize,pageIndex,filterStr三个参数
func Parse3M(r *http.Request) (filterStr string, pageIndex, pageSize int, err error) {
	pageIndex, pageSize, err = Parse2M(r)
	if err != nil {
		return "", 0, 0, err
	}

	filterStr = r.FormValue("filterStr")
	if filterStr != "" {
		if err := checkFilterStr(filterStr); err != nil {
			return "", 0, 0, err
		}
	}

	return filterStr, pageIndex, pageSize, nil
}

//Parse4M 读取请求参数的pageSize,pageIndex,filterStr,order四个参数
func Parse4M(r *http.Request) (filterStr string, pageIndex, pageSize int, order string, err error) {
	filterStr, pageIndex, pageSize, err = Parse3M(r)
	if err != nil {
		return "", 0, 0, "", err
	}

	order = r.FormValue("order")

	return filterStr, pageIndex, pageSize, order, nil
}
