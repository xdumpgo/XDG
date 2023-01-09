package main

import (
	"fmt"
	"github.com/xdumpgo/XDG/api/client"
	"github.com/xdumpgo/XDG/auth"
	"github.com/xdumpgo/XDG/cmd/qtui/uiqt"
	"github.com/xdumpgo/XDG/injection"
	"github.com/xdumpgo/XDG/modules"
	"github.com/xdumpgo/XDG/modules/dorkers"
	"github.com/xdumpgo/XDG/qtui"
	"github.com/xdumpgo/XDG/utils"
	"github.com/arl/statsviz"
	"github.com/oreans/virtualizersdk"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
)

const STEALTH_ONE_KB = 1024 / 4
const STEALTH_ONE_MB = 1024 * STEALTH_ONE_KB

const STEALTH_SIZE   = STEALTH_ONE_MB * 8

var stealth_area = [STEALTH_SIZE]uint32{0xa1a2a3a4, 0xa4a3a2a1, 0xb1a1b2a2, 0xb8a8a1a1,
	0xb6b5b3b6, 0xa2b2c2d2, 0xa9a8a2a2, 0xa0a9b9b8}

func init() {
	debug.SetMaxThreads(3000000)
	runtime.GOMAXPROCS(1024)
	modules.Scraper = &modules.ScrapeModule{
		Dorks: []string{},
		Urls: []string{},
	}
	modules.Exploiter = &modules.ExploiterModule{
		Injectables: map[string]*injection.Injection{},
	}
	modules.Dumper = &modules.DumpModule{}
	modules.AutoSheller = &modules.AutoShellModule{}
	//manager.PManager = &manager.ProxyManager{}
	/*manager.PManager = &manager.ProxyManager{
		Proxies: make([]*manager.Proxy,0),
		Client: http.Client{Timeout: time.Duration(20) * time.Second},
	}
	manager.PManager.Client.Transport = manager.CreateProxyTransport()*/

	modules.Generator = modules.NewGenerator()

	modules.AntiPublic = modules.NewAntiPublic()

	modules.Analyzer = modules.NewAnaylzer()

	utils.LogInfo("Setting up Dorkers")
	dorkers.SetupDorkers()
}

func main() {
	statsviz.RegisterDefault()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	virtualizersdk.Macro(virtualizersdk.SHARK_WHITE_START)

	if err := client.ConnectToAPIServer(); err == nil {
		if auth.Heartbeat() {
			uiqt.StartUI()
		} else {
			os.Exit(0)
		}
	} else {
		qtui.SimpleMB(nil, "Failed to connect to licensing API", "Fatal Error").Show()
	}

	if stealth_area[0] == 0x11111111 {
		fmt.Println(stealth_area[0])
	}
	virtualizersdk.Macro(virtualizersdk.SHARK_WHITE_END)
}
