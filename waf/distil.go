package waf

import "net/http"

func Distil(headers http.Header, str string) bool {
	if HasAny(str, []string{"cdn.distilnetworks.com/images/anomaly.detected.png", "distilCaptchaForm", "distilCallbackGuard"}) {
		return true
	}
	return false
}