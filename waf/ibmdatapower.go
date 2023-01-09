package waf

import (
	"net/http"
	"regexp"
)

func IBMDataPower (headers http.Header, str string) bool {
	if match1, _ := regexp.MatchString(`(OK|FAIL)`, headers.Get("X-Backside-Transport")); match1 {
		return true
	}
	return false
}