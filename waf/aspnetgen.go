package waf

import (
	"net/http"
	"regexp"
)

func ASPNetGeneric(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`iis (\d+.)+?detailed error`, str)
	if match1 || HasAny(str, []string{"potentially dangerous request querystring", "application error from being viewed remotely (for security reasons)?", "An application error occurred on the server"}) {
		return true
	}
	return false
}