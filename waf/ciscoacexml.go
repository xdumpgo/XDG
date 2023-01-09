package waf

import (
	"net/http"
	"strings"
)

func CiscoAceXML(headers http.Header, str string) bool {
	if strings.Contains(headers.Get("Server"), "ACE XML Gateway") {
		return true
	}
	return false
}