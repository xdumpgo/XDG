package waf

import (
	"net/http"
	"regexp"
)

func CloudBric(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`<title>Cloudbric.{0,5}?ERROR!`, str)
	match2, _ := regexp.MatchString(`malformed request syntax.{0,4}?invalid request message framing.{0,4}?or deceptive request routing`, str)
	if match1 || HasAny(str, []string{"Your request was blocked by Cloudbric", "please contact Cloudbric Support", "cloudbric.zendesk.com", "Cloudbric Help Center"}) || match2 {
		return true
	}
	return false
}