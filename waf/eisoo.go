package waf

import (
	"net/http"
	"strings"
	"regexp"
)

func Eisoo (headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`EisooWAF(\-AZURE)?/?`, headers.Get("Server"))
	match2, _ := regexp.MatchString(`<link.{0,10}href=\"/eisoo\-firewall\-block\.css`, str)
	match3, _ := regexp.MatchString(`&copy; \d{4} Eisoo Inc`, str)
	if match1 || match2 || strings.Contains(str, "www.eisoo.com") || match3 {
		return true
	}
	return false
}