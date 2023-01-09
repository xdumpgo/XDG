package waf

import (
	"net/http"
)

func Frontdoor (headers http.Header, str string) bool {
	if _, match1 := headers["X-Azure-Ref"]; match1 {
		return true
	}
	return false
}