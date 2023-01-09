package waf

import (
	"net/http"
	"strings"
)

func DOSArrest(headers http.Header, str string) bool {
	_, match1 := headers["X-DIS-Request-ID"]
	if match1 || strings.Contains(headers.Get("Server"), "DOSarrest") {
		return true
	}
	return false
}