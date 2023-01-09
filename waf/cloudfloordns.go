package waf

import (
	"net/http"
	"strings"
	"regexp"
)

func CloudfloorDNS(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`CloudfloorDNS(.WAF)?`, headers.Get("Server"))
	match2, _ := regexp.MatchString(`<(title|h\d{1})>CloudfloorDNS.{0,6}?Web Application Firewall Error`, str)
	if match1 || match2 || strings.Contains(str, "www.cloudfloordns.com/contact") {
		return true
	}
	return false
}