package waf

import (
	"net/http"
	"strings"
	"regexp"
)

func Anyu(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`anyu.{0,10}?the green channel`, str)
	if match1 || strings.Contains(str, "your access has been intercepted by anyu") {
		return true
	}
	return false
}