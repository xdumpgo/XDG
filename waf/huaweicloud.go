package waf

import (
	"net/http"
	"regexp"
	"strings"
)

func HuaweiCloud (headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`^HWWAFSESID=`, headers.Get("Set-Cookie"))
	if match1 || strings.Contains(headers.Get("Server"), "HuaweiCloudWAF") || HasAny(str, []string{"hwcloud.com", "hws_security@"}) {
		return true
	}
	return false
}