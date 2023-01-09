package waf

import (
	"net/http"
	"strings"
)

func Comodo(headers http.Header, str string) bool {
	if strings.Contains(headers.Get("Server"), "Protected by COMODO WAF") {
		return true
	}
	return false
}