package waf

import (
	"net/http"
	"strings"
)

func Barikode(headers http.Header, str string) bool {
	if strings.Contains(str, "<strong>barikode</strong>") {
		return true
	}
	return false
}