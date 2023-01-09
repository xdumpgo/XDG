package waf

import (
	"net/http"
	"regexp"
)

func Baidu(headers http.Header, str string) bool {
	if match1, _ := regexp.MatchString("Yunjiasu(.+)?", headers.Get("Server")); match1 {
		return true
	}
	return false
}