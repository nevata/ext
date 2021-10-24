package ext

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type User struct {
	Model
	Name string
	Age  uint8
	Addr string
}

func TestGet(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Error(err)
		return
	}

	db.AutoMigrate(&DevOrder{})
	db.AutoMigrate(&DevWhere{})
	db.AutoMigrate(&User{})

	r := gin.Default()
	r.GET("/api/user.fetch", func(c *gin.Context) {
		where, values, order, page, err := Parse(c.Request)
		if err != nil {
			HandleExcept(c.Writer, err)
			return
		}

		t.Log("where", where, "order", order, "values", values, "page", page)

		tx := db.Begin()
		defer tx.Rollback()

		if where != nil {
			query, vals, err := WhereBuild(tx, *where, values)
			if err != nil {
				HandleExcept(c.Writer, err)
				return
			}
			t.Log("query", query)
			if query != "" {
				tx = tx.Where("1=1 "+query, vals...)
			}
		}

		if order != nil {
			s := OrderByBuild(tx, *order)
			t.Log("order", s)
			if s != "" {
				tx = tx.Order(s)
			}
		}

		var total int64
		if err := tx.Debug().Model(&User{}).Count(&total).Error; err != nil {
			HandleExcept(c.Writer, err)
			return
		}

		var users []User
		if err := tx.Debug().Offset(page.Index).Limit(page.Size).
			Find(&users).Error; err != nil {
			HandleExcept(c.Writer, err)
			return
		}

		HandleSuccess(c.Writer, map[string]interface{}{
			"items": users,
			"total": total,
		})
	})

	go r.Run(":80")

	args := url.Values{}
	//分页
	args.Add("page_index", "1") //第几页，从1开始
	args.Add("page_size", "50") //每页最大条数
	//排序
	args.Add("ext_order", "1") //排序sql
	//过滤
	args.Add("ext_where", "1") //过滤sql
	args.Add("name", "rice")   //过滤条件等于
	args.Add("age", "36")      //过滤条件大于
	args.Add("addr", "普陀%")    //过滤条件LIKE

	resp, err := http.Get("http://127.0.0.1/api/user.fetch?" + args.Encode())
	if err != nil {
		t.Error(err)
		return
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(buf))
}
