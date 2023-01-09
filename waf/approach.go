package waf

import (
	"net/http"
	"regexp"
)

func Approach(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`approach.{0,10}?web application (firewall|filtering)`, str)
	match2, _ := regexp.MatchString(`approach.{0,10}?infrastructure team`, str)
	if match1 || match2 {
		return true
	}
	return false
}