package waf

import (
	"net/http"
	"strings"
)

func Alree(headers http.Header, str string) bool {
	if strings.Contains(str, "airee.cloud") || strings.Contains(headers.Get("Server"), "Airee") || strings.Contains(headers.Get("X-Cache"), "airee.cloud") {
		return true
	}
	return false
}