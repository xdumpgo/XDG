package waf

import (
	"net/http"
	"strings"
	"regexp"
)

func DynamicWeb(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`by dynamic check(.{0,10}?module)?`, str)
	if strings.Contains(headers.Get("X-403-Status-By"), "dw.inj.check") || match1 {
		return true
	}
	return false
}