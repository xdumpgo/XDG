package waf

import (
	"net/http"
	"strings"
)

func CrawlProtect(headers http.Header, str string) bool {
	if strings.Contains(headers.Get("Set-Cookie"), "crawlprotecttag=") || HasAny(str, []string{"<title>crawlprotect", "this site is protected by crawlprotect"}) {
		return true
	}
	return false
}