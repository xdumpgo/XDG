package waf

import (
	"net/http"
	"regexp"
)

func Azion(headers http.Header, str string) bool {
	if match1, _ := regexp.MatchString("Azion([-_]CDN)?", headers.Get("Server")); match1 {
		return true
	}
	return false
}