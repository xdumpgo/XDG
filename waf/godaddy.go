package waf

import (
	"net/http"
	"regexp"
)

func GoDaddy (headers http.Header, str string) bool {
	if match1, _ := regexp.MatchString(`GoDaddy (security|website firewall)`, str); match1 {
		return true
	}
	return false
}