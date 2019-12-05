package ext

import "github.com/nevata/session"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	Auth        bool
	HandlerFunc session.HandlerFunc
}
