package waf

import (
	"net/http"
	"strings"
)

func Cdnns(headers http.Header, str string) bool {
	if strings.Contains(str, "cdnsswaf application gateway") {
		return true
	}
	return false
}