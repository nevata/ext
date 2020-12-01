package ext

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//checkKeywords 过滤关键字
func checkKeywords(buf []byte) bool {
	return bytes.Contains(buf, []byte{'-', '-'}) ||
		bytes.Contains(buf, []byte{'/', '*'}) ||
		bytes.Contains(buf, []byte{'D', 'E', 'L', 'E', 'T', 'E', ' '}) ||
		bytes.Contains(buf, []byte{'I', 'N', 'S', 'E', 'R', 'T', ' '}) ||
		bytes.Contains(buf, []byte{'U', 'P', 'D', 'A', 'T', 'E', ' '}) ||
		bytes.Contains(buf, []byte{'C', 'R', 'E', 'A', 'T', 'E', ' '}) ||
		bytes.Contains(buf, []byte{'D', 'R', 'O', 'P', ' '}) ||
		bytes.Contains(buf, []byte{'T', 'R', 'U', 'N', 'C', 'A', 'T', 'E', ' '}) ||
		bytes.Contains(buf, []byte{'A', 'L', 'T', 'E', 'R', ' '}) ||
		bytes.Contains(buf, []byte{'U', 'N', 'I', 'O', 'N', ' '}) ||
		bytes.Contains(buf, []byte{'I', 'N', 'T', 'O', ' '}) ||
		bytes.Contains(buf, []byte{'J', 'O', 'I', 'N', ' '}) ||
		bytes.Contains(buf, []byte{'S', 'L', 'E', 'E', 'P', ' '})
}

//checkWhere 校验请求的过滤参数
func checkWhere(s string) bool {
	buf := []byte(strings.ToUpper(s))
	bufLen := len(buf)

	//前缀、后缀检查
	prefix1 := []byte{'A', 'N', 'D', '('}
	prefix2 := []byte{'A', 'N', 'D', ' ', '('}
	suffix := []byte{')'}

	if !bytes.HasPrefix(buf, prefix1) &&
		!bytes.HasPrefix(buf, prefix2) &&
		!bytes.HasSuffix(buf, suffix) {
		log.Println("prefix & suffix check failed: ", s)
		return false
	}

	//最一个小括号与第一个小括号成对闭合
	var q []byte
	for k, v := range buf {
		if v == '(' {
			q = append(q, v)
		} else if v == ')' {
			q = q[:len(q)-1]
		}
		if k >= 4 && k != bufLen-1 && len(q) == 0 {
			log.Println("'(' check failed: ", s)
			return false
		}
	}

	//关键字检查
	if checkKeywords(buf) {
		log.Println("keywords check failed: ", s)
		return false
	}

	return true
}

//checkOrder 检验请求的排序参数
func checkOrder(s string) bool {
	buf := []byte(strings.ToUpper(s))
	// bufLen := len(buf)

	//前缀检查
	prefix := []byte{','}

	if !bytes.HasPrefix(buf, prefix) {
		log.Println("prefix check failed: ", s)
		return false
	}

	//关键字检查
	if checkKeywords(buf) {
		log.Println("keywords check failed: ", s)
		return false
	}

	return true
}

//Parse2M 读取请求参数的pageSize,pageIndex两个参数
func Parse2M(r *http.Request) (pageIndex, pageSize int, err error) {
	if err := r.ParseForm(); err != nil {
		return 0, 0, err
	}

	pageSizeStr := r.FormValue("page_size")
	if pageSizeStr == "" {
		pageSizeStr = r.FormValue("pageSize")
		if pageSizeStr == "" {
			pageSizeStr = r.FormValue("PageSize")
		}
	}
	delete(r.Form, "page_size")
	delete(r.Form, "pageSize")
	delete(r.Form, "PageSize")

	pageSize, err = strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}

	if pageSize > 100 {
		pageSize = 100
	}

	pageIndexStr := r.FormValue("page_index")
	if pageIndexStr == "" {
		pageIndexStr = r.FormValue("pageIndex")
		if pageIndexStr == "" {
			pageIndexStr = r.FormValue("PageIndex")
		}
	}
	delete(r.Form, "page_index")
	delete(r.Form, "pageIndex")
	delete(r.Form, "PageIndex")

	pageIndex, err = strconv.Atoi(pageIndexStr)
	if err != nil {
		pageIndex = 0
	}

	return pageIndex, pageSize, nil
}

//Parse3M 读取请求参数的pageSize,pageIndex,filterStr三个参数
func Parse3M(r *http.Request) (
	filterStr string,
	pageIndex, pageSize int,
	err error,
) {
	pageIndex, pageSize, err = Parse2M(r)
	if err != nil {
		return "", 0, 0, err
	}

	filterStr = r.FormValue("filter_str")
	if filterStr == "" {
		filterStr = r.FormValue("filterStr")
		if filterStr == "" {
			filterStr = r.FormValue("FilterStr")
		}
	}
	delete(r.Form, "filter_str")
	delete(r.Form, "filterStr")
	delete(r.Form, "FilterStr")

	if filterStr != "" {
		if !checkWhere(filterStr) {
			return "", 0, 0, errors.New("[filter_str]参数不正确")
		}
	}

	return filterStr, pageIndex, pageSize, nil
}

//Parse4M 读取请求参数的pageSize,pageIndex,filterStr,orderStr四个参数
func Parse4M(r *http.Request) (
	filterStr string,
	pageIndex, pageSize int,
	orderStr string,
	err error,
) {
	filterStr, pageIndex, pageSize, err = Parse3M(r)
	if err != nil {
		return "", 0, 0, "", err
	}

	orderStr = r.FormValue("order_str")
	if orderStr != "" {
		if !checkOrder(orderStr) {
			return "", 0, 0, "", errors.New("[order_str]参数不正确")
		}
	}

	return filterStr, pageIndex, pageSize, orderStr, nil
}

//Page 数据分页结构
type Page struct {
	Index int //第几页，索引从零开始
	Size  int //每一页大小
}

//Parse 解析过滤条件、分页参数、排序字段
func Parse(r *http.Request) (
	where *uint,
	values map[string]string,
	order *uint,
	page *Page,
	err error,
) {
	pageIndex, pageSize, err := Parse2M(r)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	page = &Page{
		Index: pageIndex - 1,
		Size:  pageSize,
	}

	extWhere := r.FormValue("ext_where")
	if extWhere != "" {
		id, err := strconv.Atoi(extWhere)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		whereid := uint(id)
		where = &whereid
	}
	delete(r.Form, "ext_where")

	extOrder := r.FormValue("ext_order")
	if extOrder != "" {
		id, err := strconv.Atoi(extOrder)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		orderid := uint(id)
		order = &orderid
	}
	delete(r.Form, "ext_order")

	values = map[string]string{}
	if where != nil {
		for key, value := range r.Form {
			if len(value) > 0 && value[0] != "" {
				values[key] = value[0]
			}
		}
	}

	return
}
