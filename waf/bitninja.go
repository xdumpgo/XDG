package waf

import (
	"net/http"
)

func BitNinja(headers http.Header, str string) bool {
	if HasAny(str, []string{"Security check by BitNinja", "Visitor anti-robot validation"}) {
		return true
	}
	return false
}