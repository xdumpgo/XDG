package waf

import (
	"net/http"
	"regexp"
)

func Incapsula (headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`^incap_es.*?=`, headers.Get("Set-Cookie"))
	match2, _ := regexp.MatchString(`^visid_incap.*?=`, headers.Get("Set-Cookie"))
	if match1 || match2 || HasAny(str, []string{"incapsula incident id", "powered by incapsula", "/_Incapsula_Resource"}) {
		return true
	}
	return false
}