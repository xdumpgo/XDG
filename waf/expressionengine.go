package waf

import (
	"net/http"
	"strings"
	"regexp"
)

func ExpressionEngine (headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`^exp_track.+?=`, headers.Get("Set-Cookie"))
	match2, _ := regexp.MatchString(`^exp_last_.+?=`, headers.Get("Set-Cookie"))
	if match1 || match2 || strings.Contains(str, "invalid get data") {
		return true
	}
	return false
}