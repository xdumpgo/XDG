package uiqt

import (
	"encoding/json"
	"fmt"
	"github.com/xdumpgo/XDG/api/client"
	protocol "github.com/xdumpgo/XDG/apiproto"
	"github.com/xdumpgo/XDG/auth"
	"github.com/xdumpgo/XDG/injection"
	"github.com/xdumpgo/XDG/modules"
	"github.com/xdumpgo/XDG/qtui"
	"github.com/xdumpgo/XDG/utils"
	"github.com/SilverCory/discordrpc-go"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)


func StartUI() {
	qtui.Application = widgets.NewQApplication(len(os.Args), os.Args)
	utils.LogInfo("Created application")

	qtui.Application.SetStyle2("Fusion")
	p := qtui.Application.Palette(nil)
	p.SetColor(gui.QPalette__All, gui.QPalette__Window, gui.QColor_FromRgb2(53, 53, 53, 255))
	p.SetColor(gui.QPalette__All, gui.QPalette__WindowText, gui.QColor_FromRgb2(255,255,255,255))
	p.SetColor(gui.QPalette__All, gui.QPalette__Text, gui.QColor_FromRgb2(255,255,255,255))
	p.SetColor(gui.QPalette__All, gui.QPalette__Base, gui.QColor_FromRgb2(25,25,25,255))
	p.SetColor(gui.QPalette__All, gui.QPalette__AlternateBase, gui.QColor_FromRgb2(25,25,25,255))
	p.SetColor(gui.QPalette__All, gui.QPalette__Button, gui.QColor_FromRgb2(53, 53, 53, 255))
	p.SetColor(gui.QPalette__All, gui.QPalette__Highlight, gui.QColor_FromRgb2(142, 45, 197, 255))
	p.SetColor(gui.QPalette__All, gui.QPalette__ButtonText, gui.QColor_FromRgb2(255, 255, 255, 255))
	p.SetColor(gui.QPalette__All, gui.QPalette__BrightText, gui.QColor_FromRgb2(255,255,255,255))

	qtui.Application.SetPalette(p, "")
	utils.LogInfo("Applied custom visual styles")

	NewSettingsWindow()
	utils.LogInfo("Created Settings Widget")
	NewWhitelistWindow()
	utils.LogInfo("Created Whitelist Widget")
	NewParametersWindow()
	utils.LogInfo("Created Parameters Widget")
	NewAuthWindow()
	utils.LogInfo("Created Authentication Widget")
	NewMain()
	utils.LogInfo("Created Main Widget")
	NewSingleSiteAnalyzer()
	utils.LogInfo("Created Single Site Widget")


	utils.LogInfo("Setup widgets")

	lcdPalette := qtui.Application.Palette(nil)

	lcdPalette.SetColor(gui.QPalette__All, gui.QPalette__Light, gui.QColor_FromRgb2(255,255,255,255))
	qtui.Main.AntipubLinkCount.SetPalette(lcdPalette)
	qtui.Main.AntipubDomainCount.SetPalette(lcdPalette)
	qtui.Main.AntipubPublic.SetPalette(lcdPalette)
	qtui.Main.AntipubLoaded.SetPalette(lcdPalette)
	qtui.Main.AntipubPrivate.SetPalette(lcdPalette)
	qtui.Main.AntipubPrivateRatio.SetPalette(lcdPalette)

	qtui.Main.AntipubLinkCount.Display2(modules.AntiPublic.Count("urls"))
	qtui.Main.AntipubDomainCount.Display2(modules.AntiPublic.Count("domains"))

	qtui.Main.AntipubSizeOnDisk.SetText(fmt.Sprintf("Size On Disk: %s", utils.ByteCountIEC(modules.AntiPublic.Size())))

	qtui.Main.CustomerChatbox.AddItem("Connected.")
	go func() {
		for {
			select {
			case cmd := <-client.XDGAPI.Incoming():
				qtui.Main.CustomerChatbox.AddItem(fmt.Sprintf("[%s] %s: %s", time.Now().Format("15:04:05"), cmd.Name, cmd.Message))
				qtui.Main.CustomerChatbox.ScrollToBottom()
			case cmd := <-client.XDGAPI.UserList():
				for i:= qtui.Main.CustomerList.Count(); i >= 0; i-- {
					qtui.Main.CustomerList.TakeItem(i)
				}
				qtui.Main.CustomerList.AddItems(cmd.Users)
			case cmd := <- client.XDGAPI.Terminate():
				qtui.SimpleMB(qtui.Main, cmd.Reason, "Your session has been terminated.").Show()
				os.Exit(1)
			case cmd := <- client.XDGAPI.News():
				for _, news := range cmd.News {
					qtui.Main.NewsBox.AddItem(fmt.Sprintf("%s | %s", news.Time.Format(time.Stamp), news.Message))
				}
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
					Urls:          len(modules.Scraper.Urls),
					Injectables:   len(modules.Exploiter.Injectables),
					Rows:          modules.Dumper.Rows,
				})
			}
		}
	}()

	utils.LogInfo("Setup Discord RPC")

	presence := &discordrpc.CommandRichPresenceMessage{
		CommandMessage: discordrpc.CommandMessage{Command: "SET_ACTIVITY"},
		Args: &discordrpc.RichPresenceMessageArgs{
			Pid: os.Getpid(),
			Activity: &discordrpc.Activity{
				Details:  "XDumpGO",
				State:    "Dumping High Quality Data",
				Assets: &discordrpc.Assets{
					LargeText:    "Dumpy dump",
					LargeImageID: "xdg",
					SmallText:    "Dank Memes",
					SmallImageID: "quartz",
				},
			},
		},
	}

	presence.SetNonce()
	data, err := json.Marshal(presence)

	if err != nil {
		log.Fatal(err.Error())
	}

	_ = utils.DiscordRP.Write(string(data))

	if utils.FileExists(".xdg") {
		b, err := ioutil.ReadFile(".xdg")
		str := strings.Split(utils.Decrypt(string(b), "xNz#'%/2n4SZsB>m"), ":")
		if err == nil {
			resp := auth.Login(str[0], str[1])
			if resp.Status == "success" {
				injection.Init()
				qtui.Main.Show()
			} else {
				qtui.AuthWindow.Show()
			}
		}
	} else {
		qtui.AuthWindow.Show()
	}
	widgets.QApplication_Exec()
}