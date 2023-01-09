package waf

import (
	"net/http"
)

func BulletProof(headers http.Header, str string) bool {
	if HasAny(str, []string{"bpsMessage", "403 Forbidden Error Page", "If you arrived here due to a search"}) {
		return true
	}
	return false
}