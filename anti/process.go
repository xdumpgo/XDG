// +build windows

package anti

import (
	"bytes"
	"github.com/StackExchange/wmi"
)

var debugBlacklist = [...]string{
	"NETSTAT",
	"FILEMON",
	"PROCMON",
	"REGMON",
	"CAIN",
	"NETMON",
	"Tcpview",
	"vpcmap",
	"vmsrvc",
	"vmusrvc",
	"wireshark",
	"VBoxTray",
	"VBoxService",
	"IDA",
	"WPE PRO",
	"The Wireshark Network Analyzer",
	"WinDbg",
	"OllyDbg",
	"Colasoft Capsa",
	"Microsoft Network Monitor",
	"Fiddler",
	"SmartSniff",
	"Immunity Debugger",
	"Process Explorer",
	"PE Tools",
	"AQtime",
	"DS-5 Debug",
	"Dbxtool",
	"Topaz",
	"FusionDebug",
	"NetBeans",
	"Rational Purify",
	".NET Reflector",
	"Cheat Engine",
	"Sigma Engine",
	"codecracker",
	"x32dbg",
	"x64dbg",
	"ida",
	"charles",
	"dnspy",
	"simpleassembly",
	"peek",
	"httpanalyzer",
	"httpdebug",
	"fiddler",
	"wireshark",
	"proxifier",
	"mitmproxy",
	"ethereal",
	"airsnare",
	"smsniff",
	"smartsniff",
	"netmon",
	"processhacker",
	"killswitch",
	"codecracker",
	"ghidra",
	"Burpsuite",
	"Ghidra",
	"dnSpy",
	"Fiddler",
}

type Win32_Process struct {
	Name           string
	ExecutablePath *string
}

func checkForProc(proc string) bool {
	var dst []Win32_Process
	q := wmi.CreateQuery(&dst, "")
	err := wmi.Query(q, &dst)
	if err != nil {
		return false
	}
	for _, v := range dst {
		if bytes.Contains([]byte(v.Name), []byte(proc)) {
			return true
		}
	}
	return false
}
