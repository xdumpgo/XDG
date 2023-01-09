package waf

import (
	"net/http"
	"regexp"
	"strings"
)

func Airlock(headers http.Header, str string) bool {
	match, _ := regexp.MatchString(`^al[_-]?(sess|lb)=`, headers.Get("Set-Cookie"))
	if match || strings.Contains(str, "server detected a syntax error in your request") {
		return true
	}
	
	return false
}