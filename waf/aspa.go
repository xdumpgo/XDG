package waf

import (
	"net/http"
	"regexp"
)

func Aspa(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`ASPA[\-_]?WAF`, headers.Get("Server"))
	if _, ok := headers["ASPA-Cache-Status"]; ok || match1 {
		return true
	}
	return false
}