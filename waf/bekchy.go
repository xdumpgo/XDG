package waf

import (
	"net/http"
	"strings"
	"regexp"
)

func Bekchy(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`Bekchy.{0,10}?Access Denied.`, str)
	if match1 || strings.Contains(str, "bekchy.com/report") {
		return true
	}
	return false
}