package waf

import (
	"net/http"
	"strings"
)

func Fortiweb (headers http.Header, str string) bool {
	if (strings.Contains(headers.Get("Set-Cookie"), "FORTIWAFSID=") || strings.Contains(str, ".fgd_icon")) || (strings.Contains(str, "fgd_icon") && strings.Contains(str, "web.page.blocked") && strings.Contains(str, "url") && strings.Contains(str, "attack.id") && strings.Contains(str, "message.id") && strings.Contains(str, "client.ip")) {
		return true
	}
	return false
}