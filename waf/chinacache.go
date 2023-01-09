package waf

import (
	"net/http"
)

func ChinaCache(headers http.Header, str string) bool {
	if _, ok := headers["Powered-By-ChinaCache"]; ok {
		return true
	}
	return false
}