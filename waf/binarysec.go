package waf

import (
	"net/http"
	"strings"
)

func BinarySec(headers http.Header, str string) bool {
	_, match1 := headers["x-binarysec-via"]
	_, match2 := headers["x-binarysec-nocache"]
	if strings.Contains(headers.Get("Server"), "BinarySec") || match1 || match2 {
		return true
	}
	return false
}