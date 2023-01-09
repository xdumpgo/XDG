package waf

import (
	"net/http"
	"regexp"
)

func AlertLogic(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`<(title|h\d{1})>requested url cannot be found`, str)
	match2, _ := regexp.MatchString(`we are sorry.{0,10}?but the page you are looking for cannot be found`, str)
	if match1 || match2 || HasAny(str, []string{"back to previous page", "proceed to homepage", "reference id"}) {
		return true
	}
	return false
}