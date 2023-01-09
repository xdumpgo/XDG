package waf

import (
	"net/http"
	"strings"
)

func AnquanBao(headers http.Header, str string) bool {
	if _, ok := headers["X-Powered-By-Anquanbao"]; ok || strings.Contains(str, "aqb_cc/error/") {
		return true
	}
	return false
}