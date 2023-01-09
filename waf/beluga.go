package waf

import (
	"net/http"
	"strings"
	"regexp"
)

func Beluga(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString("^beluga_request_trail=", headers.Get("Set-Cookie"))
	if match1 || strings.Contains(headers.Get("Server"), "Beluga") {
		return true
	}
	return false
}