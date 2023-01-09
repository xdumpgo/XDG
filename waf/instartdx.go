package waf

import (
	"net/http"
	"strings"
)

func InstartDX (headers http.Header, str string) bool {
	_, match1 := headers["X-Instart-Request-ID"]
	_, match2 := headers["X-Instart-Cache"]
	_, match3 := headers["X-Instart-WL"]
	
	if (match1 || match2 || match3) || (strings.Contains(str, "the requested url was rejected") && strings.Contains(str, "please consult with your administrator") && strings.Contains(str, "your support id is")) {
		return true
	}
	return false
}