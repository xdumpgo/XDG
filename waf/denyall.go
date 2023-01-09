package waf

import (
	"net/http"
	"strings"
)

func DenyAll(headers http.Header, str string) bool {
	if strings.Contains(str, "Condition Intercepted") {
		return true
	}
	return false
}