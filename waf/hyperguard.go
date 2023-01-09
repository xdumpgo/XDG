package waf

import (
	"net/http"
	"regexp"
)

func HyperGuard (headers http.Header, str string) bool {
	if match1, _ := regexp.MatchString(`^WODSESSION=`, headers.Get("Set-Cookie")); match1 {
		return true
	}
	return false
}