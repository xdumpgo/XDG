package waf

import (
	"net/http"
	"regexp"
	"strings"
)

func GreyWizard (headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`<(title|h\d{1})>Grey Wizard`, str)
	if strings.Contains(headers.Get("Server"), "greywizard") || match1 || strings.Contains(str, "contact the website owner or Grey Wizard") || strings.Contains(str, "We've detected attempted attack or non standard traffic from your ip address") {
		return true
	}
	return false
}