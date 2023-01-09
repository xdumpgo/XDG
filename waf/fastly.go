package waf

import (
	"net/http"
)

func Fastly (headers http.Header, str string) bool {
	if _, match1 := headers["X-Fastly-Request-ID"]; match1 {
		return true
	}
	return false
}