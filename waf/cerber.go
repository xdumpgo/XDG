package waf

import (
	"net/http"
	"regexp"
)

func Cerber(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`We.re sorry.{0,10}?you are not allowed to proceed`, str)
	if match1 || HasAny(str, []string{"you request looks suspicious or similar to automated", "our server stopped processing your request", "requests from spam posting software", "<title>403 Access Forbidden"}) {
		return true
	}
	return false
}