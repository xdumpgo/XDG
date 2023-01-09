package waf

import (
	"net/http"
	"regexp"
)

func Aliyundun(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`error(s)?\.aliyun(dun)?\.(com|net)?`, str)
	match2, _ := regexp.MatchString(`cdn\.aliyun(cs)?\.com`, str)
	match3, _ := regexp.MatchString(`^aliyungf_tc=`, headers.Get("Set-Cookie"))
	if match1 || match2 || match3 {
		return true
	}
	return false
}