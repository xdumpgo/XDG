package waf

import (
	"net/http"
)

func ArvanCloud(headers http.Header, str string) bool {
	if headers.Get("Server") == "ArvanCloud" {
		return true
	}
	return false
}