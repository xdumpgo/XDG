package waf

import (
	"net/http"
	"strings"
	"regexp"
)

func F5BigIPAPM (headers http.Header, str string) bool {
	match1, _ := regexp.MatchString(`^LastMRH_Session`, headers.Get("Set-Cookie"))
	match2, _ := regexp.MatchString(`^MRHSession`, headers.Get("Set-Cookie"))
	
	match3, _ := regexp.MatchString(`Big([-_])?IP`, headers.Get("Server"))
	
	match4, _ := regexp.MatchString(`^F5_fullWT`, headers.Get("Set-Cookie"))
	match5, _ := regexp.MatchString(`^F5_HT_shrinked`, headers.Get("Set-Cookie"))
	
	if (match1 && match2) || (match2 && match3) || (match4 || match5) {
		return true
	}
	return false
}

func F5BigIPASM (headers http.Header, str string) bool {
	if strings.Contains(str, "the requested url was rejected") && strings.Contains(str, "please consult with your administrator") {
		return true
	}
	return false
}

func F5BigIPLTM (headers http.Header, str string) bool {
	if strings.Contains(headers.Get("Set-Cookie"), "bigipserver") || strings.Contains(headers.Get("X-Cnection"), "close") {
		return true
	}
	return false
}

func F5FirePass (headers http.Header, str string) bool {
	if (strings.Contains(headers.Get("Location"), "my.logon.php3") && strings.Contains(headers.Get("Set-Cookie"), "VHOST")) || (strings.Contains(headers.Get("Set-Cookie"), "F5_fire") && strings.Contains(headers.Get("Set-Cookie"), "F5_passid_shrinked")) {
		return true
	}
	return false
}

func F5TrafficShield (headers http.Header, str string) bool {
	if strings.Contains(headers.Get("Set-Cookie"), "ASINFO=") || strings.Contains(headers.Get("Server"), "F5-TrafficShield") {
		return true
	}
	return false
}