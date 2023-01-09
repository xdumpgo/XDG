package waf

import (
	"net/http"
	"regexp"
)

func EdgeCast(headers http.Header, str string) bool {
	if match1, _ := regexp.MatchString(`^EC(D|S)(.*)?`, headers.Get("Server")); match1 {
		return true
	}
	return false
}