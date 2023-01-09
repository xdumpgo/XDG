package waf

import (
	"net/http"
	"regexp"
)

func Imunify360 (headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`imunify360.{0,10}`, headers.Get("Server"))
	match2, _ := regexp.MatchString(`protected.by.{0,10}?imunify360`, str)
	match3, _ := regexp.MatchString(`powered.by.{0,10}?imunify360`, str)
	match4, _ := regexp.MatchString(`imunify360.preloader`, str)
	if match1 || match2 || match3 || match4 {
		return true
	}
	return false
}