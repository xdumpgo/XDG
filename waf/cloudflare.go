package waf

import (
	"net/http"
	"strings"
	"regexp"
)

func Cloudflare(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString("cloudflare", headers.Get("Server"))
	match2, _ := regexp.MatchString("cloudflare[-_]nginx", headers.Get("Server"))
	_, match3 := headers["cf-ray"]
	if match1 || match2 || match3 || strings.Contains(headers.Get("Set-Cookie"), "__cfduid") {
		return true
	}
	return false
}