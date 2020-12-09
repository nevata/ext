package ext

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//Model 基本模型的定义
type Model struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt JSONTime  `json:"create_at"`
	UpdatedAt JSONTime  `json:"update_at"`
	DeletedAt *JSONTime `sql:"index" json:"-"`
}

//DevOrder 排序字段
type DevOrder struct {
	ID      uint   `gorm:"primary_key"`
	Columns string //字段
	Desc    string //备注
}

//OrderByBuild 根据ID返回排序字段
func OrderByBuild(db *sql.DB, id uint) string {
	row := db.QueryRow("SELECT columns FROM dev_orders WHERE id = ?", id)
	var orderby string
	err := row.Scan(&orderby)
	if err == sql.ErrNoRows {
		return ""
	}
	if err != nil {
		PrintErr(err)
		return ""
	}
	return orderby
}

//DevWhere SQL过滤条件
type DevWhere struct {
	ID         uint   `gorm:"primary_key"`
	Conditions string //条件
	Desc       string //备注
}

//WhereBuild 根据ID返回过滤条件
func WhereBuild(db *sql.DB, id uint, keys map[string]string) (string, []interface{}, error) {
	row := db.QueryRow("SELECT conditions FROM dev_wheres WHERE id = ?", id)
	var s string
	err := row.Scan(&s)
	if err == sql.ErrNoRows {
		return "", nil, nil
	}
	if err != nil {
		return "", nil, err
	}

	type ValueType uint8

	const (
		VTInt      ValueType = iota + 1 //1 int
		VTFloat                         //2 float
		VTString                        //3 string
		VTDate                          //4 date
		VTTime                          //5 time
		VTDateTime                      //6 datetime
		VTStringEx                      //7 %string%
	)

	var conditions []struct {
		Key      string    `json:"key"`
		Val      string    `json:"val"`
		VT       ValueType `json:"vt"`
		Required bool      `json:"required"`
	}

	err = json.NewDecoder(strings.NewReader(s)).Decode(&conditions)
	if err != nil {
		return "", nil, err
	}

	where := ""
	vals := []interface{}{}
	for _, condition := range conditions {
		val, ok := keys[condition.Key]
		if condition.Required || ok {
			if where != "" {
				where += " AND "
			}
			where += condition.Val
			if ok {
				switch condition.VT {
				case VTInt:
					v, e := strconv.ParseInt(val, 10, 64)
					if e != nil {
						return "", nil, fmt.Errorf("参数[%s]类型不正确", condition.Key)
					}
					vals = append(vals, sql.Named(condition.Key, v))
				case VTFloat:
					v, e := strconv.ParseFloat(val, 64)
					if e != nil {
						return "", nil, fmt.Errorf("参数[%s]类型不正确", condition.Key)
					}
					vals = append(vals, sql.Named(condition.Key, v))
				case VTDate:
					v, e := time.Parse(dateLayout, val)
					if e != nil {
						return "", nil, fmt.Errorf("参数[%s]类型不正确", condition.Key)
					}
					vals = append(vals, sql.Named(condition.Key, v))
				case VTTime:
					v, e := time.Parse("15:04:05", val)
					if e != nil {
						return "", nil, fmt.Errorf("参数[%s]类型不正确", condition.Key)
					}
					vals = append(vals, sql.Named(condition.Key, v))
				case VTDateTime:
					v, e := time.Parse(timeLayout, val)
					if e != nil {
						return "", nil, fmt.Errorf("参数[%s]类型不正确", condition.Key)
					}
					vals = append(vals, sql.Named(condition.Key, v))
				case VTStringEx:
					vals = append(vals, sql.Named(condition.Key, "%"+val+"%"))
				default:
					vals = append(vals, sql.Named(condition.Key, val))
				}
			}
		}
	}

	return where, vals, nil
}
