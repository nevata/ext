package ext

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
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
func OrderByBuild(db *gorm.DB, id uint) string {
	row := db.Raw("SELECT columns FROM dev_orders WHERE id = ?", id)
	var orderby string
	err := row.Scan(&orderby).Error
	if err == gorm.ErrRecordNotFound {
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

//ValueType 值类型
type ValueType uint8

//支持的值类型
const (
	VTInt      ValueType = iota + 1 //1 int
	VTFloat                         //2 float
	VTString                        //3 string
	VTDate                          //4 date
	VTTime                          //5 time
	VTDateTime                      //6 datetime
	VTStringEx                      //7 %string%
	VTSubQuery                      //8
)

//Condition 条件结构定义
type Condition struct {
	Key         string       `json:"key"`
	Val         string       `json:"val"`
	VT          ValueType    `json:"vt"`
	KeyRequired bool         `json:"key_required"` //必传参数
	ValRequired bool         `json:"val_required"` //必加条件
	Conditions  []*Condition `json:"conditions"`
}

//
func encodeJSON(v interface{}) string {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func parseCondition(
	condition *Condition,
	keys map[string]string,
	val string,
	builder *strings.Builder,
	vals *[]interface{},
) error {
	switch condition.VT {
	case VTInt:
		v, e := strconv.ParseInt(val, 10, 64)
		if e != nil {
			return fmt.Errorf("参数[%s]类型不正确", condition.Key)
		}
		*vals = append(*vals, sql.Named(condition.Key, v))
	case VTFloat:
		v, e := strconv.ParseFloat(val, 64)
		if e != nil {
			return fmt.Errorf("参数[%s]类型不正确", condition.Key)
		}
		*vals = append(*vals, sql.Named(condition.Key, v))
	case VTDate:
		v, e := time.Parse(dateLayout, val)
		if e != nil {
			return fmt.Errorf("参数[%s]类型不正确", condition.Key)
		}
		*vals = append(*vals, sql.Named(condition.Key, v))
	case VTTime:
		v, e := time.Parse("15:04:05", val)
		if e != nil {
			return fmt.Errorf("参数[%s]类型不正确", condition.Key)
		}
		*vals = append(*vals, sql.Named(condition.Key, v))
	case VTDateTime:
		v, e := time.Parse(timeLayout, val)
		if e != nil {
			return fmt.Errorf("参数[%s]类型不正确", condition.Key)
		}
		*vals = append(*vals, sql.Named(condition.Key, v))
	case VTStringEx:
		*vals = append(*vals, sql.Named(condition.Key, "%"+val+"%"))
	case VTSubQuery:
		return parseConditions(condition.Conditions, keys, builder, vals)
	default:
		*vals = append(*vals, sql.Named(condition.Key, val))
	}

	return nil
}

func parseConditions(
	conditions []*Condition,
	keys map[string]string,
	builder *strings.Builder,
	vals *[]interface{},
) error {
	for _, condition := range conditions {
		val, ok := keys[condition.Key]
		if condition.KeyRequired && !ok {
			return fmt.Errorf("[%s]必传参数", condition.Key)
		}
		if condition.ValRequired || ok {
			if builder.Len() > 0 {
				builder.WriteString(" AND ")
			}
			builder.WriteString(condition.Val)
			if ok {
				if err := parseCondition(
					condition, keys, val, builder, vals); err != nil {
					return err
				}
				log.Println("parseConditions:", vals)
			}
		}
	}
	return nil
}

//WhereBuild 根据ID返回过滤条件
func WhereBuild(
	db *gorm.DB,
	id uint,
	keys map[string]string,
) (
	string,
	[]interface{},
	error,
) {
	row := db.Raw("SELECT conditions FROM dev_wheres WHERE id = ?", id)
	var s string
	err := row.Scan(&s).Error
	if err == gorm.ErrRecordNotFound {
		return "", nil, nil
	}
	if err != nil {
		return "", nil, err
	}

	var conditions []*Condition
	err = json.NewDecoder(strings.NewReader(s)).Decode(&conditions)
	if err != nil {
		return "", nil, err
	}

	builder := strings.Builder{}
	vals := []interface{}{}
	if err := parseConditions(conditions, keys, &builder, &vals); err != nil {
		return "", nil, err
	}

	return builder.String(), vals, nil
}
