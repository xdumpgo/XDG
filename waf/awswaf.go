package waf

import (
	"net/http"
	"regexp"
)

func AWSWaf(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString("^aws.?alb=", headers.Get("Set-Cookie"))
	match2, _ := regexp.MatchString("aws.?elb", headers.Get("Server"))
	_, match3 := headers["X-AMZ-ID"]
	_, match4 := headers["X-AMZ-Request-ID"]
	if match1 || match2 || match3 || match4 {
		return true
	}
	return false
}