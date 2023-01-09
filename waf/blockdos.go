package waf

import (
	"net/http"
	"strings"
)

func BlockDOS(headers http.Header, str string) bool {
	if strings.Contains(headers.Get("Server"), "blockdos.net") {
		return true
	}
	return false
}