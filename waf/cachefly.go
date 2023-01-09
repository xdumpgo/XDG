package waf

import (
	"net/http"
	"strings"
	"regexp"
)

func CacheFly(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`^cfly_req.*=`, headers.Get("Set-Cookie"))
	if strings.Contains(headers.Get("BestCDN"), "Cachefly") || match1 {
		return true
	}
	return false
}