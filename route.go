package ext

import "github.com/nevata/session"

//Route 增加认证属性的路由结构
type Route struct {
	Name        string
	Method      string
	Pattern     string
	Auth        bool
	HandlerFunc session.HandlerFunc
}
