package waf

import (
	"net/http"
	"strings"
)

const (
	AESECURE = iota
	ALREE
	AIRLOCK
	ALERTLOGIC
	ALIYUNDUN
	ANQUANBAO
	ANYU
	APPROACH
	ARMOR
	ARVANCLOUD
	ASPA
	ASPNETGEN
	ASTAR
	AWSWAF
	AZION
	BAIDU
	BARIKODE
	BARRACUDA
	BEKCHY
	BELUGA
	BINARYSEC
	BITNINJA
	BLOCKDOS
	BLUEDON
	BULLETPROOF
	CACHEFLY
	
)

type WAF struct {
	Test func(headers http.Header, str string) bool
	Tampers []string
}

var Wafs []*WAF

func init() {
	Wafs = append(Wafs, &WAF{AESecure, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Airlock, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{AlertLogic, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Aliyundun, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Alree, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{AnquanBao, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Anyu, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Approach, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Armor, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{ArvanCloud, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Aspa, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{ASPNetGeneric, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Astra, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{AWSWaf, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Azion, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Baidu, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Barikode, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Barracuda, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Bekchy, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Beluga, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{BinarySec, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{BitNinja, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{BlockDOS, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{BlueDon, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{BulletProof, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{CacheFly, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{CacheWall, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Cdnns, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Cerber, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{ChinaCache, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{ChuangYu, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{CiscoAceXML, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{CloudBric, []string{"luanginx", "space2comment"}})
	Wafs = append(Wafs, &WAF{Cloudflare, []string{"luanginx", "space2comment"}})
	Wafs = append(Wafs, &WAF{CloudfloorDNS, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{CloudFront, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Comodo, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{CrawlProtect, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{DenyAll, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Distil, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{DOSArrest, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{DotDefender, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{DynamicWeb, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{EdgeCast, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Eisoo, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{ExpressionEngine, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{F5BigIPAPM, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{F5BigIPASM, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{F5BigIPLTM, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{F5FirePass, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{F5TrafficShield, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Fastly, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Fortiweb, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Frontdoor, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{GoDaddy, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{GreyWizard, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{HuaweiCloud, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{HyperGuard, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{IBMDataPower, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Imunify360, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{Incapsula, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{InstartDX, []string{"space2comment"}})
	Wafs = append(Wafs, &WAF{ModSecurity, []string{"modsec", "between", "space2comment"}})
}

func IsWAF(resp *http.Response, body string) (bool, []string) {
	for _,waf := range Wafs {
		if waf.Test(resp.Header, body) {
			return true, waf.Tampers
		}
	}
	return true, []string{"space2comment"}
}

func HasAny(thing string, shit []string) bool {
	for _, v := range shit {
		if strings.Contains(thing, v) {
			return true
		}
	}
	return false
}