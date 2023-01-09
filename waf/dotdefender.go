package waf

import "net/http"

func DotDefender(headers http.Header, str string) bool {
	_, match1 := headers["X-dotDefender-denied"]
	if match1 || HasAny(str, []string{"dotdefender blocked your request", "Applicure is the leading provider of web application security"}) {
		return true
	}
	return false
}