package waf

import (
	"net/http"
	"strings"
)

func Armor(headers http.Header, str string) bool {
	if strings.Contains(str, "blocked by website protection from armor") || strings.Contains(str, "please create an armor support ticket") {
		return true
	}
	return false
}