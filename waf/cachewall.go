package waf

import (
	"net/http"
	"regexp"
)

func CacheWall(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString("Varnish", headers.Get("Server"))
	_, match2 := headers["X-Varnish"]
	_, match3 := headers["X-Cachewall-Action"]
	_, match4 := headers["X-Cachewall-Reasons"]
	match5, _ := regexp.MatchString(`403 naughty.{0,10}?not nice!`, str)
	if match1 || match2 || match3 || match4 || match5 || HasAny(str, []string{"security by cachewall", "varnish cache server"}) {
		return true
	}
	return false
}