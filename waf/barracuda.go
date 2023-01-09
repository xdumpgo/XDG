package waf

import (
	"net/http"
	"strings"
	"regexp"
)

func Barracuda(headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`^barra_counter_session=`, headers.Get("Set-Cookie"))
	match2, _ := regexp.MatchString(`^BNI__BARRACUDA_LB_COOKIE=`, headers.Get("Set-Cookie"))
	match3, _ := regexp.MatchString(`^BNI_persistence=`, headers.Get("Set-Cookie"))
	match4, _ := regexp.MatchString(`^BN[IE]S_.*?=`, headers.Get("Set-Cookie"))
	if match1 || match2 || match3 || match4 || strings.Contains(str, "Barracuda.Networks") {
		return true
	}
	return false
}