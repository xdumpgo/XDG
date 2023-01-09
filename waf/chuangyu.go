package waf

import (
	"net/http"
)

func ChuangYu(headers http.Header, str string) bool {
	if HasAny(str, []string{"www.365cyd.com", "help.365cyd.com/cyd-error-help.html?code=403"}) {
		return true
	}
	return false
}