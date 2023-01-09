package waf

import (
	"net/http"
	"strings"
)

func BlueDon(headers http.Header, str string) bool {
	if strings.Contains(headers.Get("Server"), "BDWAF") || strings.Contains(str, "bluedon web application firewall") {
		return true
	}
	return false
}