package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/xdumpgo/XDG/modules"
	"github.com/xdumpgo/XDG/modules/dorkers"
	"github.com/xdumpgo/XDG/utils"
	"github.com/logrusorgru/aurora"
	"gopkg.in/ini.v1"
)

func init() {
	utils.LogInfo(fmt.Sprintf("[%s] Startup", aurora.Magenta("XDG")))
	_, _ = utils.SetConsoleTitle("XDumpGO | Created by Zertex#0001 & xiQQ#6974")

	modules.Scraper = &modules.ScrapeModule{}
	modules.Exploiter = &modules.ExploiterModule{}
	modules.Dumper = &modules.DumpModule{}

	utils.LogInfo("Checking files")
	if !utils.FileExists("proxies.txt") {
		file, _ := os.Create("proxies.txt")
		utils.LogInfo("Creating urls.txt")
		file.Close()
	}
	if !utils.FileExists("dorks.txt") {
		file, _ := os.Create("dorks.txt")
		utils.LogInfo("Creating dorks.txt")
		file.Close()
	}
	if !utils.FileExists("urls.txt") {
		utils.LogInfo("Creating urls.txt")
		file, _ := os.Create("urls.txt")
		file.Close()
	}
	if !utils.FileExists("injectables.txt") {
		file, _ := os.Create("injectables.txt")
		utils.LogInfo("Creating injectables.txt")
		file.Close()
	}
	var err error
	utils.LogInfo("Loading config")
	if utils.CFG, err = ini.Load("config.ini"); err != nil {
		utils.CFG, _ = ini.Load([]byte(`
[main]
Threads         = 50
Pages           = 2
Timeouts        = 5
BatchMode       = false
DumpAll         = false
AutoThreads     = true
TechError       = true
TechUnion       = true
TechBlind       = false
Level           = 0
TargetedDump    = false

[cluster]
Enabled    = false
Beacon     = 0.0.0.0
Port       = 3579
SmallStats = false

[dorkers]
Google      = false
Bing        = false
AOL         = false
MyWebSearch = false
Yahoo       = false
DuckDuckGo  = false
Ecosia      = false
Qwant       = false
StartPage   = false
Yandex      = false
`))
		utils.CFG.SaveTo("config.ini")
	}

	utils.LogInfo("Setting up Dorkers")
	dorkers.SetupDorkers()

	utils.LogInfo("Setting up UI components")
	//ui.CreateUI()

	utils.LogInfo("Done!")

	time.Sleep(750 * time.Millisecond)
	utils.CallClear()
}

func main() {
	if utils.FileExists(".xdg") {
		b, err := ioutil.ReadFile(".xdg")
		str := strings.Split(utils.Decrypt(string(b), "xNz#'%/2n4SZsB>m"), ":")
		if err == nil {
			ts := &auth.ToServer{
				Username:   str[0],
				Password:   str[1],
				PacketType: "authenticate",
			}
			auth.ClientUsername = str[0]
			auth.Password = str[1]
			utils.LogInfo("Authenticating...")
			resp := auth.DoTransaction(ts)
			if resp.Status == "success" {
				auth.Expiry = resp.Expiry
				ui.InfoText.SetText(fmt.Sprintf("Welcome, %s, your key expires on: [cyan]%s", auth.ClientUsername, auth.Expiry.Format("Mon Jan 2 15:04:05 2006")))
				injection.Init()
				ui.Pages.SwitchToPage("Main")
			}
		}
	}

	if err := ui.StartUI(); err != nil {
		panic(err)
	}
}
