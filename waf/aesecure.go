package waf

import (
	"net/http"
	"strings"
)

func AESecure(headers http.Header, str string) bool {
	_, match := headers["aeSecure-code"]
	if match || strings.Contains(str, "aesecure_denied.png") {
		return true
	}
	return false
}