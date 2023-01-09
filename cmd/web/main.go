package main

import (
	"fmt"
	"github.com/xdumpgo/XDG/api/client"
	protocol "github.com/xdumpgo/XDG/apiproto"
	"github.com/xdumpgo/XDG/auth"
	"github.com/xdumpgo/XDG/injection"
	"github.com/xdumpgo/XDG/manager"
	"github.com/xdumpgo/XDG/modules"
	"github.com/xdumpgo/XDG/modules/dorkers"
	"github.com/xdumpgo/XDG/utils"
	"github.com/xdumpgo/XDG/web"
	"github.com/webview/webview"
	"os"
	"time"
)

func init() {
	modules.Scraper = &modules.ScrapeModule{
		Dorks: []string{},
		Urls: []string{},
	}
	modules.Exploiter = &modules.ExploiterModule{
		Injectables: map[string]*injection.Injection{},
	}
	modules.Dumper = &modules.DumpModule{}
	modules.AutoSheller = &modules.AutoShellModule{}
	manager.PManager = &manager.ProxyManager{}

	modules.Generator = modules.NewGenerator()

	utils.LogInfo("Setting up Dorkers")
	dorkers.SetupDorkers()
}

func main() {
	if err := client.ConnectToAPIServer(); err == nil {
		if auth.Heartbeat() {
			go func() {
				for {
					select {
					case cmd := <-client.XDGAPI.Incoming():
						cmd.Name = cmd.Name
					case cmd := <-client.XDGAPI.UserList():
						cmd.Users = cmd.Users
					case cmd := <- client.XDGAPI.Terminate():
						cmd.Reason = cmd.Reason
						os.Exit(1)
					case <- client.XDGAPI.StatusUpdate():
						var index int
						var end int
						switch utils.Module {
						case "Scraper":
							index = modules.Scraper.Index
							end = len(modules.Scraper.Dorks)
							break
						case "Exploiter":
							index = modules.Exploiter.Index
							end = len(modules.Scraper.Urls)
							break
						case "Dumper":
							index = modules.Dumper.Index
							end = len(modules.Exploiter.Injectables)
							break
						case "AutoSheller":
							index = modules.AutoSheller.Index
							end = len(modules.Exploiter.Injectables)
						}


						client.XDGAPI.Send(protocol.StatsUpdate{
							CurrentModule: utils.Module,
							Runtime:       time.Since(utils.StartTime),
							Index:         index,
							End:           end,
							Threads:       len(utils.GlobalSem),
							Workers:	   len(utils.WorkerSem),
							Urls:          len(modules.Scraper.Urls),
							Injectables:   len(modules.Exploiter.Injectables),
							Rows:          modules.Dumper.Rows,
						})
					}
				}
			}()
			go web.StartWebUI(":1337")

			w := webview.New(false)
			defer w.Destroy()
			w.SetTitle(fmt.Sprintf("XDumpGO v%s - Created by Zertex#0001 & xiQQ#0001", auth.Version))
			w.SetSize(1920, 1080, webview.HintNone)
			w.Navigate("http://localhost:1337")
			w.Run()
		} else {
			os.Exit(0)
		}
	} else {
		utils.LogError( "Failed to connect to licensing API")
	}
}