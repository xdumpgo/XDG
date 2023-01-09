package waf

import (
	"net/http"
	"regexp"
)

func Astra(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`^cz_astra_csrf_cookie`, headers.Get("Set-Cookie"))
	if match1 || HasAny(str, []string{"astrawebsecurity.freshdesk.com", "www.getatra.com/assets/images"}) {
		return true
	}
	return false
}