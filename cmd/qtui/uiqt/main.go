package uiqt

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"github.com/xdumpgo/XDG/api/client"
	"github.com/xdumpgo/XDG/auth"
	"github.com/xdumpgo/XDG/injection"
	"github.com/xdumpgo/XDG/manager"
	"github.com/xdumpgo/XDG/modules"
	"github.com/xdumpgo/XDG/modules/dorkers"
	"github.com/xdumpgo/XDG/qtui"
	"github.com/xdumpgo/XDG/utils"
	"github.com/xdumpgo/XDG/web"
	"github.com/pkg/browser"
	"github.com/spf13/viper"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)


func StopModules(bool) {
	if utils.Module != "Idle" {
		if utils.Module == "Stopping" {
			manager.PManager.CancelAll()
			if !utils.IsClosed(utils.Kill) {
				close(utils.Kill)
			}
		} else {
			utils.Module = "Stopping"
			if !utils.IsClosed(utils.Done) {
				close(utils.Done)
			}
		}
	} else {
		qtui.SimpleMB(qtui.Main, "Already stopping!", "Error")
	}
}

func NewMain() *qtui.MainWindow {
	qtui.Main = qtui.NewMainWindow(nil)

	qtui.Main.SetWindowTitle(fmt.Sprintf("XDumpGO v%s | Created by Quartz Inc. | https://invite.gg/quartzinc", auth.Version))

	qtui.Main.ExploiterInjectablesTable.HorizontalHeader().ResizeSection(0, 160)

	//qtui.Main.DumperTableWidget.HorizontalHeader().ResizeSection(0, 160)

	qtui.Main.TabControl.ConnectCurrentChanged(func(index int) {
		if index == 6 {
			qtui.Main.CleanerUrlsTextBox.SetPlainText(strings.Join(modules.Scraper.Urls, "\r\n"))
		}
	})

	qtui.Main.NewsBox.SetWordWrap(true)

	qtui.Main.CustomerChatbox.SetWordWrap(true)

	qtui.Main.CleanerBtn.ConnectClicked(func(bool) {
		loading := qtui.NewLoadingWindow(qtui.Main)
		loading.ProgressBar.SetMaximum(len(modules.Scraper.Urls))
		loading.ProgressBar.SetValue(0)
		loading.Label.SetText("Cleaning your urls, please wait...")
		go loading.Show()

		cleanRg := regexp.MustCompile(`^http.+\?.+=`)

		go func() {
			modules.Scraper.Urls = []string{}
			domains := make(map[string]interface{})
			i := 0
			scanner := bufio.NewScanner(strings.NewReader(qtui.Main.CleanerUrlsTextBox.ToPlainText()))
			for scanner.Scan() {
				i++
				loading.ProgressBar.SetValue(i)
				u, err := url.Parse(scanner.Text())
				if err != nil {
					continue
				}
				if qtui.Main.CleanerDuplicateDomains.IsChecked() {
					if _, ok := domains[u.Hostname()]; ok {
						continue
					}
					domains[u.Hostname()] = nil
				}

				if qtui.Main.CleanerQueryParam.IsChecked() {
					if !cleanRg.MatchString(u.String()) {
						continue
					}
				}

				modules.Scraper.Urls = append(modules.Scraper.Urls, u.String())
			}
			qtui.Main.CleanerUrlsTextBox.SetPlainText(strings.Join(modules.Scraper.Urls, "\r\n"))
			os.Mkdir("cleaned-urls", 0755)
			f := utils.CreateFileTimeStamped("output", "cleaned-urls")
			f.WriteString(strings.Join(modules.Scraper.Urls, "\r\n"))
			f.Close()
			loading.Hide()
			loading.Close()
		}()
	})

	qtui.Main.WebUILaunchBtn.ConnectClicked(func(bool) {
		go web.StartWebUI(qtui.Main.WebUIAddress.Text())
		var address string
		if len(strings.Split(qtui.Main.WebUIAddress.Text(), ":")[0]) == 0 {
			address = "localhost"+qtui.Main.WebUIAddress.Text()
		} else {
			address = qtui.Main.WebUIAddress.Text()
		}
		browser.OpenURL("http://"+address)
	})

	go func() {
		for {
			select {
			case <- time.After(time.Second):
				qtui.Main.StatsModule.SetText(utils.Module)
				qtui.Main.StatsCurrentTimeLcd.Display(time.Now().Format("15:04:05"))
				if utils.Module != "Idle" {
					qtui.Main.StatsRuntimeLcd.Display(utils.FmtDuration(time.Since(utils.StartTime)))
				}
				qtui.Main.DataDorksLcd.Display2(len(modules.Scraper.Dorks))
				qtui.Main.DataUrlsLcd.Display2(len(modules.Scraper.Urls))
				qtui.Main.DataInjectablesLcd.Display2(len(modules.Exploiter.Injectables))
				qtui.Main.DataInjectablesLcd.Repaint()
				qtui.Main.DataProxiesLcd.Display2(len(manager.PManager.Proxies))
				qtui.Main.StatsRequestsLcd.Display2(utils.RequestCounter)
				qtui.Main.StatsRPSLcd.Display2(int(utils.RateCounter.Rate()))
				qtui.Main.StatsErrorLcd.Display2(utils.ErrorCounter)
				qtui.Main.DataThreadsLcd.Display2(len(utils.GlobalSem))
				qtui.Main.DataWorkersLcd.Display2(len(utils.WorkerSem))

				qtui.Main.DumpStatsTablesLcd.Display2(modules.Dumper.Tables)
				qtui.Main.DumpStatsColumnsLcd.Display2(modules.Dumper.Columns)
				qtui.Main.DumpStatsRowsLcd.Display2(modules.Dumper.Rows)
				qtui.Main.DumpStatsRpm.Display(utils.RPMCounter.String())

				switch utils.Module {
				case "Scraper":
					qtui.Main.ParserProgress.SetMaximum(len(modules.Scraper.Dorks))
					qtui.Main.ParserProgress.SetValue(modules.Scraper.Index)
					//qtui.Main.ParserProgress.SetFormat(fmt.Sprintf("%.2f%% %d/%d", float64(modules.Scraper.Index/len(modules.Scraper.Dorks)), modules.Scraper.Index, len(modules.Scraper.Dorks)))
				case "Exploiter":
					qtui.Main.ExploiterProgress.SetMaximum(len(modules.Scraper.Urls))
					qtui.Main.ExploiterProgress.SetValue(modules.Exploiter.Index)
					//qtui.Main.ExploiterProgress.SetFormat(fmt.Sprintf("%.2f%% %d/%d", float64(modules.Exploiter.Index/len(modules.Scraper.Urls)), modules.Exploiter.Index, len(modules.Scraper.Urls)))
				case "Dumper":
					qtui.Main.DumperProgress.SetMaximum(len(modules.Exploiter.Injectables))
					qtui.Main.DumperProgress.SetValue(modules.Dumper.Index)
					//qtui.Main.DumperProgress.SetFormat(fmt.Sprintf("%.2f%% %d/%d", float64(modules.Dumper.Index/len(modules.Exploiter.Injectables)), modules.Dumper.Index, len(modules.Exploiter.Injectables)))
				case "AutoSheller":
					qtui.Main.AutoShellProgress.SetMaximum(len(modules.Exploiter.Injectables))
					qtui.Main.AutoShellProgress.SetValue(modules.AutoSheller.Index)
					//qtui.Main.AutoShellProgress.SetFormat(fmt.Sprintf("%.2f%% %d/%d", float64(modules.AutoSheller.Index/len(modules.Exploiter.Injectables)), modules.AutoSheller.Index, len(modules.Exploiter.Injectables)))
				case "Generator":
					qtui.Main.GeneratorProgress.SetMaximum(len(modules.Generator.Parameters))
					qtui.Main.GeneratorProgress.SetValue(modules.Generator.Index)
					//qtui.Main.GeneratorProgress.SetFormat(fmt.Sprintf("%.2f%% %d/%d", float64(modules.Generator.Index/len(modules.Generator.Parameters)), modules.Generator.Index, len(modules.Generator.Parameters)))
				case "AntiPublic":
					qtui.Main.AntipubProgress.SetMaximum(len(modules.Scraper.Urls))
					qtui.Main.AntipubProgress.SetValue(modules.AntiPublic.Index)
				}
			}
		}
	}()

	/*--------- GENERATOR ----------*/

	var paramNames []string
	for pI, p := range modules.Generator.Parameters {
		paramNames = append(paramNames, fmt.Sprintf("%s-%s", p.Name,pI))
	}
	fmt.Println(paramNames)
	qtui.Main.ComboBoxParameters.AddItems(paramNames)

	// Handles reading the data from the files and showing it
	qtui.Main.ComboBoxParameters.ConnectCurrentIndexChanged(func(index int) {
		if index := qtui.Main.ComboBoxParameters.CurrentIndex(); index > -1 {
			prefix := strings.Split(qtui.Main.ComboBoxParameters.CurrentText(), "-")[1]

			var data []string
			if data = modules.Generator.Parameters[prefix].GetData(); len(data) == 0 {
				if f, err := os.Open(modules.Generator.Parameters[prefix].FilePath); err == nil {
					scanner := bufio.NewScanner(f)
					for scanner.Scan() {
						data = append(data, scanner.Text())
					}
				}
			}

			qtui.Main.PlainTextEditParameters.SetPlainText(strings.Join(data, "\n"))
		}
	})
	qtui.Main.ComboBoxParameters.SetCurrentIndex(0)
	fmt.Println(qtui.Main.ComboBoxParameters.CurrentText())
	qtui.Main.PlainTextEditParameters.SetPlainText(strings.Join(modules.Generator.Parameters[strings.Split(qtui.Main.ComboBoxParameters.CurrentText(), "-")[1]].GetData(), "\n"))

	qtui.Main.PlainTextEditPatterns.SetPlainText(strings.Join(viper.GetStringSlice("generator.Patterns"), "\n"))

	//Handles saving the data written on the box to the text files
	qtui.Main.PlainTextEditParameters.ConnectFocusOutEvent(func(event *gui.QFocusEvent) {
		if index := qtui.Main.ComboBoxParameters.CurrentIndex(); index > -1 {
			prefix := strings.Split(qtui.Main.ComboBoxParameters.CurrentText(), "-")[1]

			modules.Generator.Parameters[prefix].Data = []string{}
			scanner := bufio.NewScanner(strings.NewReader(qtui.Main.PlainTextEditParameters.ToPlainText()))
			for scanner.Scan() {
				modules.Generator.Parameters[prefix].Data = append(modules.Generator.Parameters[prefix].Data, scanner.Text())
			}
			ioutil.WriteFile(modules.Generator.Parameters[prefix].FilePath, []byte(strings.Join(modules.Generator.Parameters[prefix].Data, "\r\n")), os.ModePerm)
		}
	})

	qtui.Main.PlainTextEditPatterns.ConnectFocusOutEvent(func(event *gui.QFocusEvent) {
		modules.Generator.Patterns = []*modules.Pattern{}
		scanner := bufio.NewScanner(strings.NewReader(qtui.Main.PlainTextEditPatterns.ToPlainText()))
		re, _ := regexp.Compile(`\(([^()]+)\)`)
		for scanner.Scan() {
			if strings.TrimSpace(scanner.Text()) != "" {
				if m := re.FindAllString(scanner.Text(), -1); m != nil{
					//Debug
					fmt.Println(m)
					modules.Generator.Patterns = append(modules.Generator.Patterns, &modules.Pattern{
						Prefixes:    m,
						TotalDorks: 1,
						PatternStr:  scanner.Text(),
					})
				}
			}
		}
	})

	qtui.Main.GeneratorStartBtn.ConnectClicked(func(checked bool) {
		if utils.Module != "Idle" && utils.Module != "Stopping" {
			qtui.SimpleMB(qtui.Main, fmt.Sprintf("Already running module '%s'", utils.Module), "Error")
			return
		}
		qtui.Main.TableGenerator.Clear()
		qtui.Main.TableGenerator.ClearContents()
		qtui.Main.TableGenerator.SetRowCount(0)

		go func(){
			defer func() {
				utils.Module = "Idle"
			}()
			utils.Module = "Generator"
			modules.Generator.Start()
		}()
	})

	qtui.Main.GeneratorStopBtn.ConnectClicked(StopModules)

	viper.SetDefault("generator.limit", true)
	qtui.Main.GeneratorLimiterCheckbox.SetChecked(viper.GetBool("generator.limit"))
	qtui.Main.GeneratorLimiterCheckbox.ConnectClicked(func(checked bool) {
		viper.Set("generator.limit", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("generator.max", 5000)
	qtui.Main.GeneratorLimiterSpinbox.SetMaximum(9999999999)
	qtui.Main.GeneratorLimiterSpinbox.SetValue(viper.GetInt("generator.max"))
	qtui.Main.GeneratorLimiterSpinbox.ConnectValueChanged(func(i int) {
		viper.Set("generator.max", i)
		viper.WriteConfig()
	})

	viper.SetDefault("scraper.customparam", "")
	qtui.Main.ParserCustomParams.SetText(viper.GetString("scraper.customparam"))
	qtui.Main.ParserCustomParams.ConnectTextChanged(func(text string) {
		viper.Set("scraper.customparam", text)
		viper.WriteConfig()
	})

	viper.SetDefault("scraper.engines.google", false)
	qtui.Main.ParserGoogleCheckbox.SetChecked(dorkers.Dorkers["Google"].Enabled)
	qtui.Main.ParserGoogleCheckbox.ConnectClicked(func(checked bool) {
		dorkers.Dorkers["Google"].Enabled = checked
		viper.Set("scraper.Engines.Google", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("scraper.engines.googleapi", false)
	/*qtui.Main.ParserGoogleAPICheckbox.Hide()
	//qtui.Main.ParserGoogleAPICheckbox.SetChecked(dorkers.Dorkers["GoogleAPI"].Enabled)
	qtui.Main.ParserGoogleAPICheckbox.ConnectClicked(func(checked bool) {
		dorkers.Dorkers["GoogleAPI"].Enabled = checked
		viper.Set("scraper.Engines.GoogleAPI", checked)
		viper.WriteConfig()
	})*/

	viper.SetDefault("scraper.engines.bing", false)
	qtui.Main.ParserBingCheckbox.SetChecked(dorkers.Dorkers["Bing"].Enabled)
	qtui.Main.ParserBingCheckbox.ConnectClicked(func(checked bool) {
		dorkers.Dorkers["Bing"].Enabled = checked
		viper.Set("scraper.Engines.Bing", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("scraper.engines.aol", false)
	qtui.Main.ParserAOLCheckbox.SetChecked(dorkers.Dorkers["AOL"].Enabled)
	qtui.Main.ParserAOLCheckbox.ConnectClicked(func(checked bool) {
		dorkers.Dorkers["AOL"].Enabled = checked
		viper.Set("scraper.Engines.AOL", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("scraper.engines.mywebsearch", false)
	qtui.Main.ParserMWSCheckbox.SetChecked(dorkers.Dorkers["MyWebSearch"].Enabled)
	qtui.Main.ParserMWSCheckbox.ConnectClicked(func(checked bool) {
		dorkers.Dorkers["MyWebSearch"].Enabled = checked
		viper.Set("scraper.Engines.MyWebSearch", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("scraper.engines.duckduckgo", false)
	qtui.Main.ParserDuckDuckGoCheckbox.Hide()
	qtui.Main.ParserDuckDuckGoCheckbox.SetChecked(dorkers.Dorkers["DuckDuckGo"].Enabled)
	qtui.Main.ParserDuckDuckGoCheckbox.ConnectClicked(func(checked bool) {
		dorkers.Dorkers["DuckDuckGo"].Enabled = checked
		viper.Set("scraper.Engines.DuckDuckGo", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("scraper.engines.ecosia", false)
	qtui.Main.ParserEcosiaCheckbox.SetChecked(dorkers.Dorkers["Ecosia"].Enabled)
	qtui.Main.ParserEcosiaCheckbox.ConnectClicked(func(checked bool) {
		dorkers.Dorkers["Ecosia"].Enabled = checked
		viper.Set("scraper.Engines.Ecosia", checked)
		viper.WriteConfig()
	})

	/*viper.SetDefault("scraper.engines.qwant", false)
	qtui.Main.ParserQwantCheckbox.Hide()
	qtui.Main.ParserQwantCheckbox.SetChecked(dorkers.Dorkers["Qwant"].Enabled)
	qtui.Main.ParserQwantCheckbox.ConnectClicked(func(checked bool) {
		dorkers.Dorkers["Qwant"].Enabled = checked
		viper.Set("scraper.Engines.Qwant", checked)
		viper.WriteConfig()
	})*/

	viper.SetDefault("scraper.engines.startpage", false)
	qtui.Main.ParserStartPageCheckbox.SetChecked(dorkers.Dorkers["StartPage"].Enabled)
	qtui.Main.ParserStartPageCheckbox.ConnectClicked(func(checked bool) {
		dorkers.Dorkers["StartPage"].Enabled = checked
		viper.Set("scraper.Engines.StartPage", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("scraper.engines.yahoo", false)
	qtui.Main.ParserYahooCheckbox.SetChecked(dorkers.Dorkers["Yahoo"].Enabled)
	qtui.Main.ParserYahooCheckbox.ConnectClicked(func(checked bool) {
		dorkers.Dorkers["Yahoo"].Enabled = checked
		viper.Set("scraper.Engines.Yahoo", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("scraper.engines.yandex", false)
	qtui.Main.ParserYandexCheckbox.SetChecked(dorkers.Dorkers["Yandex"].Enabled)
	qtui.Main.ParserYandexCheckbox.ConnectClicked(func(checked bool) {
		dorkers.Dorkers["Yandex"].Enabled = checked
		viper.Set("scraper.Engines.Yandex", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("scraper.filter", true)
	qtui.Main.ParserFilterUrls.SetChecked(viper.GetBool("scraper.Filter"))
	qtui.Main.ParserFilterUrls.ConnectClicked(func(checked bool) {
		viper.Set("scraper.Filter", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("scraper.pages", 2)
	qtui.Main.ParserPagesSpinbox.SetValue(viper.GetInt("scraper.Pages"))
	qtui.Main.ParserPagesSpinbox.ConnectValueChanged(func(i int) {
		viper.Set("scraper.Pages", i)
		viper.WriteConfig()
	})

	qtui.Main.ParserLoadDorksBtn.ConnectClicked(func(bool) {
		wd, _ := os.Getwd()
		ofd := widgets.NewQFileDialog2(qtui.Main, "Select your Dorks", wd, "Text / List files (*.txt *.lst)")

		ofd.ConnectFileSelected(func(file string) {
			f, err := os.Open(file)
			if err != nil {
				return
			}
			modules.Scraper.Dorks = []string{}
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				modules.Scraper.Dorks = append(modules.Scraper.Dorks, scanner.Text())
			}
			f.Close()
			qtui.Main.ParserDorksTextbox.SetPlainText(strings.Join(modules.Scraper.Dorks, "\r\n"))
		})

		ofd.Show()
		modules.Scraper.Index = 0
	})

	qtui.Main.ParserClearDorksBtn.ConnectClicked(func(bool) {
		modules.Scraper.Dorks = make([]string, 0)
		qtui.Main.ParserDorksTextbox.SetPlainText("")
		modules.Scraper.Index = 0
	})

	qtui.Main.ParserDorksTextbox.ConnectTextChanged(func() {
		modules.Scraper.Dorks = []string{}
		scanner := bufio.NewScanner(strings.NewReader(qtui.Main.ParserDorksTextbox.ToPlainText()))
		for scanner.Scan() {
			modules.Scraper.Dorks = append(modules.Scraper.Dorks, scanner.Text())
		}
	})

	qtui.Main.ParserStartBtn.ConnectClicked(func(bool) {
		if utils.Module != "Idle" && utils.Module != "Stopping" {
			go qtui.SimpleMB(qtui.Main, fmt.Sprintf("Already running module '%s'", utils.Module), "Error").Show()
			return
		}

		if len(modules.Scraper.Dorks) == 0 {
			go qtui.SimpleMB(qtui.Main, "Please load dorks.", "Error").Show()
			return
		}

		if qtui.SettingsWindow.LoadProxiesFromAPI.IsChecked() {
			resp, err := http.Get(qtui.SettingsWindow.LineEdit.Text())
			if err != nil {
				qtui.SimpleMB(qtui.Main, "Failed to get proxylist\n" + err.Error(), "Error")
				return
			}
			defer resp.Body.Close()
			raw, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				qtui.SimpleMB(qtui.Main, "Failed to get proxylist\n" + err.Error(), "Error")
				return
			}
			manager.PManager.Proxies = []*manager.Proxy{}
			manager.PManager.LoadScanner(bufio.NewScanner(bytes.NewReader(raw)), qtui.SettingsWindow.ProxyTypeComboBox.CurrentIndex())
		}

		f := sync.WaitGroup{}
		if modules.Scraper.Index > 0 {
			f.Add(1)
			q := qtui.NewYesNo(qtui.Main, "Looks like you have another scrape paused, would you like to continue this?", "Just a moment")
			q.ConnectAccepted(func() {
				f.Done()
			})
			q.ConnectRejected(func() {
				modules.Scraper.Index = 0
				f.Done()
			})
			q.Show()
		}

		go func() {
			defer func() {
				utils.Module = "Idle"
			}()
			f.Wait()
			utils.Module = "Scraper"
			modules.Scraper.Start()
		}()
	})

	qtui.Main.ParserStopBtn.ConnectClicked(StopModules)

	viper.SetDefault("exploiter.threads", 100)
	qtui.Main.ExploiterThreads.SetValue(viper.GetInt("exploiter.threads"))
	qtui.Main.ExploiterThreads.ConnectValueChanged(func(i int) {
		viper.Set("exploiter.threads", i)
		viper.WriteConfig()
	})

	viper.SetDefault("exploiter.workers", 30)
	qtui.Main.ExploiterWorkers.SetValue(viper.GetInt("exploiter.workers"))
	qtui.Main.ExploiterWorkers.ConnectValueChanged(func(i int) {
		viper.Set("exploiter.workers", i)
		viper.WriteConfig()
	})

	viper.SetDefault("exploiter.technique.error", true)
	qtui.Main.ExploiterErrorCheckbox.SetChecked(viper.GetBool("exploiter.technique.error"))
	qtui.Main.ExploiterErrorCheckbox.ConnectClicked(func(checked bool) {
		viper.Set("exploiter.technique.error", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("exploiter.technique.union", true)
	qtui.Main.ExploiterUnionCheckbox.SetChecked(viper.GetBool("exploiter.technique.union"))
	qtui.Main.ExploiterUnionCheckbox.ConnectClicked(func(checked bool) {
		viper.Set("exploiter.technique.union", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("exploiter.technique.blind", false)
	qtui.Main.ExploiterBlindCheckbox.SetChecked(viper.GetBool("exploiter.technique.blind"))
	qtui.Main.ExploiterBlindCheckbox.ConnectClicked(func(checked bool) {
		viper.Set("exploiter.technique.blind", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("exploiter.technique.stacked", false)
	qtui.Main.ExploiterStackedCheckbox.SetChecked(viper.GetBool("exploiter.Techniques.Stacked"))
	qtui.Main.ExploiterStackedCheckbox.ConnectClicked(func(checked bool) {
		viper.Set("exploiter.Technique.Stacked", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("exploiter.heuristics", true)
	qtui.Main.ExploiterHeuristsicsCheckbox.SetChecked(viper.GetBool("exploiter.Heuristics"))
	qtui.Main.ExploiterHeuristsicsCheckbox.ConnectClicked(func(checked bool) {
		viper.Set("exploiter.Heuristics", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("exploiter.intensity", 2)
	qtui.Main.ExploiterIntensityCombo.SetCurrentIndex(viper.GetInt("exploiter.Intensity"))
	qtui.Main.ExploiterIntensityCombo.ConnectCurrentIndexChanged(func(index int) {
		viper.Set("exploiter.Intensity", index)
		viper.WriteConfig()
	})

	qtui.Main.ExploiterLoadUrlsBtn.ConnectClicked(func(bool) {
		wd, _ := os.Getwd()
		ofd := widgets.NewQFileDialog2(qtui.Main, "Select your Urls", wd, "Text (*.txt)")
		ofd.Show()
		ofd.ConnectFileSelected(func(file string) {
			f, err := os.Open(file)
			if err != nil {
				return
			}
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				modules.Scraper.Urls = append(modules.Scraper.Urls, scanner.Text())
			}
			f.Close()
			//uiParseDorksTxtbx.SetPlainText(strings.Join(modules.Scraper.Urls, "\r\n"))
			modules.Exploiter.Index = 0
		})
	})

	qtui.Main.ExploiterInjectablesTable.ConnectCustomContextMenuRequested(func(pos *core.QPoint) {
		if qtui.Main.ExploiterInjectablesTable.RowCount() == 0 {
			return
		}
		table := qtui.Main.ExploiterInjectablesTable
		index := table.CurrentRow()
		i := table.Item(index, 0).Text()
		_u, _ := url.Parse(i)
		if injectable, ok := modules.Exploiter.Injectables[_u.Hostname()]; ok {
			menu := widgets.NewQMenu2(injectable.Base.String(), table)
			menu.SetWindowTitle("Actions")
			skip := menu.AddAction(fmt.Sprintf("Open Analyzer - %s", injectable.Base.String()))
			skip.ConnectTriggered(func(bool) {
				modules.Analyzer.Mount(injectable)
				qtui.SingleSiteWindow.Show()
			})
			menu.Exec2(qtui.Main.ExploiterInjectablesTable.Viewport().MapToGlobal(pos), skip)

			menu.Show()
			menu.ActivateWindow()
		}
	})

	qtui.Main.ExploiterClearUrlsBtn.ConnectClicked(func(bool) {
		modules.Scraper.Urls = []string{}
		modules.Exploiter.Index = 0
	})

	qtui.Main.ExploiterStartBtn.ConnectClicked(func(bool) {
		if utils.Module != "Idle" && utils.Module != "Stopping" {
			go qtui.SimpleMB(qtui.Main, fmt.Sprintf("Already running module '%s'", utils.Module), "Error").Show()
			return
		}
		if len(modules.Scraper.Urls) == 0 {
			go qtui.SimpleMB(qtui.Main, "Please load or scrape urls.", "Error").Show()
		} else {
			f := sync.WaitGroup{}
			if modules.Exploiter.Index > 0 {
				f.Add(1)
				q := qtui.NewYesNo(qtui.Main, "Looks like you have another scan paused, would you like to continue this?", "Just a moment")
				q.ConnectAccepted(func() {
					f.Done()
				})
				q.ConnectRejected(func() {
					modules.Exploiter.Index = 0
					qtui.Main.ExploiterInjectablesTable.ClearContents()
					qtui.Main.ExploiterInjectablesTable.SetRowCount(0)
					f.Done()
				})
				q.Show()
			}

			if qtui.SettingsWindow.LoadProxiesFromAPI.IsChecked() {
				resp, err := http.Get(qtui.SettingsWindow.LineEdit.Text())
				if err != nil {
					qtui.SimpleMB(qtui.Main, "Failed to get proxylist\n" + err.Error(), "Error")
					return
				}
				defer resp.Body.Close()
				raw, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					qtui.SimpleMB(qtui.Main, "Failed to get proxylist\n" + err.Error(), "Error")
					return
				}
				manager.PManager.Proxies = []*manager.Proxy{}
				manager.PManager.LoadScanner(bufio.NewScanner(bytes.NewReader(raw)), qtui.SettingsWindow.ProxyTypeComboBox.CurrentIndex())
			}

			go func() {
				defer func() {
					utils.Module = "Idle"
				}()
				f.Wait()
				utils.Module = "Exploiter"
				modules.Exploiter.Start(modules.Scraper.Urls)
			}()
		}
	})

	qtui.Main.ExploiterStopBtn.ConnectClicked(StopModules)

	viper.SetDefault("dumper.threads", 40)
	qtui.Main.DumperThreads.SetValue(viper.GetInt("dumper.threads"))
	qtui.Main.DumperThreads.ConnectValueChanged(func(i int) {
		viper.Set("dumper.threads", i)
		viper.WriteConfig()
	})

	viper.SetDefault("dumper.workers", 20)
	qtui.Main.DumperWorkers.SetValue(viper.GetInt("dumper.workers"))
	qtui.Main.DumperWorkers.ConnectValueChanged(func(i int) {
		viper.Set("dumper.workers", i)
		viper.WriteConfig()
	})

	viper.SetDefault("dumper.targeted", true)
	qtui.Main.DumperTargetedCheckbox.SetChecked(viper.GetBool("dumper.Targeted"))
	qtui.Main.DumperTargetedCheckbox.ConnectClicked(func(checked bool) {
		viper.Set("dumper.Targeted", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("dumper.keepblanks", false)
	qtui.Main.DumperKeepBlanksCheckbox.SetChecked(viper.GetBool("dumper.KeepBlanks"))
	qtui.Main.DumperKeepBlanksCheckbox.ConnectClicked(func(checked bool) {
		viper.Set("dumper.KeepBlanks", checked)
		viper.WriteConfig()
	})

	qtui.Main.DumperLoadInjectablesBtn.ConnectClicked(func(bool) {
		wd, _ := os.Getwd()
		ofd := widgets.NewQFileDialog2(qtui.Main, "Select your Injectables", wd, "Text (*.txt)")
		ofd.Show()
		ofd.ConnectFileSelected(func(file string) {
			f, err := os.Open(file)
			if err != nil {
				return
			}
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				inj := injection.ParseUrlString(scanner.Text())
				if inj != nil {
					modules.Exploiter.Injectables[inj.Base.Hostname()] = inj
				}
			}
			f.Close()
			//uiParseDorksTxtbx.SetPlainText(strings.Join(modules.Scraper.Urls, "\r\n"))
			modules.Dumper.Index = 0
			modules.AutoSheller.Index = 0
		})
	})

	qtui.Main.DumperClearInjectablesBtn.ConnectClicked(func(bool) {
		modules.Exploiter.Injectables = make(map[string]*injection.Injection)
		modules.Dumper.Index = 0
		modules.AutoSheller.Index = 0
	})

	qtui.Main.DumperStartBtn.ConnectClicked(func(bool) {
		if utils.Module != "Idle" && utils.Module != "Stopping" {
			go qtui.SimpleMB(qtui.Main, fmt.Sprintf("Already running module '%s'", utils.Module), "Error").Show()
			return
		}
		if len(modules.Exploiter.Injectables) == 0 {
			go qtui.SimpleMB(qtui.Main, "Please load or test injectables.", "Error").Show()
		} else {
			f := sync.WaitGroup{}
			if modules.Dumper.Index > 0 {
				f.Add(1)
				q := qtui.NewYesNo(qtui.Main, "Looks like you have another scrape paused, would you like to continue this?", "Just a moment")
				q.ConnectAccepted(func() {
					f.Done()
				})
				q.ConnectRejected(func() {
					modules.Scraper.Index = 0
					qtui.Main.DumperTableWidget.SetRowCount(0)
					f.Done()
				})
				q.Show()
			}

			if qtui.SettingsWindow.LoadProxiesFromAPI.IsChecked() {
				resp, err := http.Get(qtui.SettingsWindow.LineEdit.Text())
				if err != nil {
					qtui.SimpleMB(qtui.Main, "Failed to get proxylist\n" + err.Error(), "Error")
					return
				}
				defer resp.Body.Close()
				raw, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					qtui.SimpleMB(qtui.Main, "Failed to get proxylist\n" + err.Error(), "Error")
					return
				}
				manager.PManager.Proxies = []*manager.Proxy{}
				manager.PManager.LoadScanner(bufio.NewScanner(bytes.NewReader(raw)), qtui.SettingsWindow.ProxyTypeComboBox.CurrentIndex())
			}

			go func() {
				defer func() {
					utils.Module = "Idle"
				}()
				utils.Module = "Dumper"

				var arr []*injection.Injection
				for _, inj := range modules.Exploiter.Injectables {
					arr = append(arr, inj)
				}
				f.Wait()

				modules.Dumper.Start(arr)
			}()
		}
	})

	qtui.Main.DumperStopBtn.ConnectClicked(StopModules)

	qtui.Main.DumperOpenAnalyzer.ConnectClicked(func(bool) {
		qtui.SingleSiteWindow.Show()
	})

	qtui.Main.DumperTableWidget.ConnectCustomContextMenuRequested(func(pos *core.QPoint) {
		if qtui.Main.DumperTableWidget.CurrentRow() == 0 {
			return
		}
		table := qtui.Main.DumperTableWidget
		index := table.CurrentRow()
		if injectable, ok := modules.Dumper.UIMap[index]; ok {
			menu := widgets.NewQMenu2(injectable.Base.String(), table)
			skip := menu.AddAction("Skip")
			skip.ConnectTriggered(func(bool) {
				injectable.SkipLock.Lock()
				defer injectable.SkipLock.Unlock()
				select {
				case <-injectable.Skip:
					return
				default:
					close(injectable.Skip)
				}
			})
			menu.Exec2(qtui.Main.DumperTableWidget.Viewport().MapToGlobal(pos), skip)

			menu.Show()
			menu.ActivateWindow()
		}
	})

	viper.SetDefault("dumper.minrows", false)
	qtui.Main.DumperMinRowsCheckbox.SetChecked(viper.GetBool("dumper.MinRows"))
	qtui.Main.DumperMinRowsCheckbox.ConnectClicked(func(checked bool) {
		viper.Set("dumper.MinRows", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("dumper.minrowcount", 500)
	qtui.Main.DumperMinRowsSpinbox.SetValue(viper.GetInt("dumper.MinRowCount"))
	qtui.Main.DumperMinRowsSpinbox.ConnectValueChanged(func(i int) {
		viper.Set("dumper.MinRowCount", i)
		viper.WriteConfig()
	})

	viper.SetDefault("dumper.usedios", true)
	qtui.Main.DumperDIOSCheckbox.SetChecked(viper.GetBool("dumper.UseDIOS"))
	qtui.Main.DumperDIOSCheckbox.ConnectClicked(func(checked bool) {
		viper.Set("dumper.UseDIOS", checked)
		viper.WriteConfig()
	})

	qtui.DumperWhitelistWindow.ConnectAccepted(func() {
		viper.WriteConfig()
	})

	viper.SetDefault("dumper.autoskipval", 1000)
	qtui.Main.DumperAutoSkipSpinbox.SetValue(viper.GetInt("dumper.autoskipval"))
	qtui.Main.DumperAutoSkipSpinbox.ConnectValueChanged(func(i int) {
		viper.Set("dumper.autoskipval", i)
		viper.WriteConfig()
	})

	viper.SetDefault("dumper.autoskip", true)
	qtui.Main.DumperAutoSkip.SetChecked(viper.GetBool("dumper.autoskip"))
	qtui.Main.DumperAutoSkip.ConnectClicked(func(checked bool) {
		viper.Set("dumper.autoskip", checked)
		viper.WriteConfig()
	})

	qtui.Main.ChatSendMessage.ConnectClicked(func(bool) {
		if len(qtui.Main.ChatMessagebox.Text()) > 0 {
			err := client.XDGAPI.SendMessage(qtui.Main.ChatMessagebox.Text())
			if err != nil {
				qtui.SimpleMB(qtui.Main, err.Error(), "Error").Show()
			}
			qtui.Main.ChatMessagebox.SetText("")
		}
	})
	qtui.Main.ChatMessagebox.ConnectReturnPressed(func() {
		if len(qtui.Main.ChatMessagebox.Text()) > 0 {
			err := client.XDGAPI.SendMessage(qtui.Main.ChatMessagebox.Text())
			if err != nil {
				qtui.SimpleMB(qtui.Main, err.Error(), "Error").Show()
			}
			qtui.Main.ChatMessagebox.SetText("")
		}
	})

	viper.SetDefault("exploiter.DBMS.MySQL", true)
	qtui.Main.ExploiterMySQL.SetChecked(viper.GetBool("exploiter.DBMS.MySQL"))
	qtui.Main.ExploiterMySQL.ConnectClicked(func(checked bool) {
		viper.Set("exploiter.DBMS.MySQL", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("exploiter.DBMS.Oracle", true)
	qtui.Main.ExploiterOracle.SetChecked(viper.GetBool("exploiter.DBMS.Oracle"))
	qtui.Main.ExploiterOracle.ConnectClicked(func(checked bool) {
		viper.Set("exploiter.DBMS.Oracle", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("exploiter.DBMS.PostGreSQL", true)
	qtui.Main.ExploiterPostgreSQL.SetChecked(viper.GetBool("exploiter.DBMS.PostGreSQL"))
	qtui.Main.ExploiterPostgreSQL.ConnectClicked(func(checked bool) {
		viper.Set("exploiter.DBMS.PostGreSQL", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("exploiter.DBMS.MSSQL", true)
	qtui.Main.ExploiterMSSQL.SetChecked(viper.GetBool("exploiter.DBMS.MSSQL"))
	qtui.Main.ExploiterMSSQL.ConnectClicked(func(checked bool) {
		viper.Set("exploiter.DBMS.MSSQL", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("autosheller.ASP", false)
	qtui.Main.AutoShellASPCheckbox.Hide()
	qtui.Main.AutoShellASPCheckbox.SetChecked(viper.GetBool("autosheller.ASP"))
	qtui.Main.AutoShellASPCheckbox.ConnectClicked(func(checked bool) {
		viper.Set("autosheller.ASP", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("autosheller.PHP", true)
	qtui.Main.AutoShellPHPCheckbox.SetChecked(viper.GetBool("autosheller.PHP"))
	qtui.Main.AutoShellPHPCheckbox.ConnectClicked(func(checked bool) {
		viper.Set("autosheller.PHP", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("autosheller.key", "sdfadsgh4513sdGG435341FDGWWDFGDFHDFGDSFGDFSGDFG")
	qtui.Main.AutoShellKey.SetText(viper.GetString("autosheller.key"))
	qtui.Main.AutoShellKey.ConnectTextChanged(func(text string) {
		viper.Set("autosheller.key", text)
		viper.WriteConfig()
	})

	viper.SetDefault("autosheller.filename", "4O4.php")
	qtui.Main.AutoShellFile.SetText(viper.GetString("autosheller.filename"))
	qtui.Main.AutoShellFile.ConnectTextChanged(func(text string) {
		viper.Set("autosheller.filename", text)
		viper.WriteConfig()
	})

	qtui.Main.AutoShellLoadInjectablesBtn.ConnectClicked(func(bool) {
		wd, _ := os.Getwd()
		ofd := widgets.NewQFileDialog2(qtui.Main, "Select your Injectables", wd, "Text (*.txt)")
		ofd.Show()
		ofd.ConnectFileSelected(func(file string) {
			f, err := os.Open(file)
			if err != nil {
				return
			}
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				inj := injection.ParseUrlString(scanner.Text())
				if inj != nil {
					modules.Exploiter.Injectables[inj.Base.Hostname()] = inj
				}
			}
			f.Close()
			//uiParseDorksTxtbx.SetPlainText(strings.Join(modules.Scraper.Urls, "\r\n"))
			modules.Dumper.Index = 0
			modules.AutoSheller.Index = 0
		})
	})

	qtui.Main.AutoShellClearInjectablesBtn.ConnectClicked(func(bool) {
		modules.Exploiter.Injectables = make(map[string]*injection.Injection)
		modules.Dumper.Index = 0
		modules.AutoSheller.Index = 0
	})

	qtui.Main.AutoShellStartBtn.ConnectClicked(func(bool) {
		if utils.Module != "Idle" && utils.Module != "Stopping" {
			qtui.SimpleMB(qtui.Main, fmt.Sprintf("Already running module '%s'", utils.Module), "Error")
			return
		}
		if len(modules.Exploiter.Injectables) == 0 {
			qtui.SimpleMB(qtui.Main, "Please load or test injectables.", "Error").Show()
		} else {
			f := sync.WaitGroup{}
			if modules.Scraper.Index > 0 {
				f.Add(1)
				q := qtui.NewYesNo(qtui.Main, "Looks like you have another run paused, would you like to continue this?", "Just a moment")
				q.ConnectAccepted(func() {
					f.Done()
				})
				q.ConnectRejected(func() {
					modules.AutoSheller.Index = 0
					modules.Dumper.Index = 0
					qtui.Main.AutoShellTable.ClearContents()
					qtui.Main.AutoShellTable.SetRowCount(0)
					f.Done()
				})
				q.Show()
			}

			go func() {
				defer func() {
					utils.Module = "Idle"
				}()
				utils.Module = "AutoSheller"
				f.Wait()

				var arr []*injection.Injection
				for _, inj := range modules.Exploiter.Injectables {
					arr = append(arr, inj)
				}

				modules.AutoSheller.Start(arr)
			}()
		}
	})

	viper.SetDefault("antipub.autosync", false)
	qtui.Main.AntipubAutoSync.SetChecked(viper.GetBool("antipub.autosync"))
	qtui.Main.AntipubAutoSync.ConnectClicked(func(checked bool) {
		viper.Set("antipub.autosync", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("antipub.savepublic", false)
	qtui.Main.AntipubSavePublic.SetChecked(viper.GetBool("antipub.savepublic"))
	qtui.Main.AntipubSavePublic.ConnectClicked(func(checked bool) {
		viper.Set("antipub.savepublic", checked)
		viper.WriteConfig()
	})

	qtui.Main.AntipubClearUrls.ConnectClicked(func(bool) {
		modules.Scraper.Urls = make([]string, 0)
	})

	qtui.Main.AntipubLoadUrls.ConnectClicked(func(bool) {
		wd, _ := os.Getwd()
		ofd := widgets.NewQFileDialog2(qtui.Main, "Select your Urls", wd, "Text (*.txt)")
		ofd.Show()
		ofd.ConnectFileSelected(func(file string) {
			f, err := os.Open(file)
			if err != nil {
				return
			}
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				modules.Scraper.Urls = append(modules.Scraper.Urls, scanner.Text())
			}
			f.Close()
			//uiParseDorksTxtbx.SetPlainText(strings.Join(modules.Scraper.Urls, "\r\n"))
			modules.Exploiter.Index = 0
			qtui.Main.AntipubLoaded.Display2(len(modules.Scraper.Urls))
		})
	})

	qtui.Main.AntipubExportDB.ConnectClicked(func(bool) {
		wd, _ := os.Getwd()
		sfd := widgets.NewQFileDialog2(qtui.Main, "Select your output file", wd, "Text (*.txt)")
		sfd.SetFileMode(widgets.QFileDialog__AnyFile)
		sfd.SetAcceptMode(widgets.QFileDialog__AcceptSave)
		sfd.Show()
		sfd.ConnectFileSelected(func(file string) {
			var err error
			var rows *sql.Rows
			switch modules.AntiPublic.GetSelected() {
			case modules.DOMAINS:
				rows, err = modules.AntiPublic.DB.Query("SELECT domain FROM domains")
			case modules.URLS:
				rows, err = modules.AntiPublic.DB.Query("SELECT url FROM urls")
			}
			if err != nil {
				qtui.SimpleMB(qtui.Main, "There was an error getting data from your db", "Error").Show()
				return
			}

			f, err := os.Create(file)
			if err != nil {
				qtui.SimpleMB(qtui.Main, "There was an error creating your output file.", "Error").Show()
				return
			}
			loading := qtui.SimpleMB(qtui.Main, "Exporting your database, please wait...", "Just a moment")
			loading.Show()
			var val string
			var index int
			for rows.Next() {
				index++
				if rows.Scan(&val) != nil {
					continue
				}
				f.WriteString(val + "\r\n")
			}

			loading.SetText(fmt.Sprintf("Successfully exported %d lines", index))

		})
	})

	qtui.Main.AntipubLinkMode.ConnectClicked(func(checked bool) {
		qtui.Main.AntipubDomainMode.SetChecked(false)
	})

	qtui.Main.AntipubDomainMode.ConnectClicked(func(checked bool) {
		qtui.Main.AntipubLinkMode.SetChecked(false)
	})

	qtui.Main.AntipubLoadToDB.ConnectClicked(func(bool) {
		wd, _ := os.Getwd()
		ofd := widgets.NewQFileDialog2(qtui.Main, "Select your list", wd, "Text (*.txt)")
		ofd.Show()
		ofd.ConnectFileSelected(func(file string) {
			i, err := modules.AntiPublic.LoadToDB(file)
			if err != nil {
				fmt.Println(err.Error())
				go qtui.SimpleMB(qtui.Main, "An error was encountered: " + err.Error(), "Error")
				return
			}
			qtui.SimpleMB(qtui.Main, fmt.Sprintf("Added %d rows to the database", i), "Loaded list to DB").Show()
			qtui.Main.AntipubDomainCount.Display2(modules.AntiPublic.Count("domains"))
			qtui.Main.AntipubLinkCount.Display2(modules.AntiPublic.Count("urls"))
			qtui.Main.AntipubSizeOnDisk.SetText(fmt.Sprintf("Size On Disk: %s", utils.ByteCountIEC(modules.AntiPublic.Size())))
		})
	})


	qtui.Main.AntipubStart.ConnectClicked(func(bool) {
		if utils.Module != "Idle" && utils.Module != "Stopping" {
			go qtui.SimpleMB(qtui.Main, fmt.Sprintf("Already running module '%s'", utils.Module), "Error").Show()
			return
		}

		if len(modules.Scraper.Urls) == 0 {
			go qtui.SimpleMB(qtui.Main, "Please load or scrape urls.", "Error").Show()
			return
		}

		go func() {
			defer func() {
				utils.Module = "Idle"
			}()
			utils.Module = "AntiPublic"
			modules.AntiPublic.Start(modules.Scraper.Urls)
		}()
	})

	qtui.Main.AntipubStop.ConnectClicked(StopModules)

	qtui.Main.AutoShellStopBtn.ConnectClicked(StopModules)

	qtui.Main.ActionExit.ConnectTriggered(func(bool) {
		qtui.Application.CloseAllWindowsDefault()
	})
	
	qtui.Main.ActionOpen_Settings.ConnectTriggered(func(bool) {
		qtui.SettingsWindow.Show()
	})

	qtui.Main.DumperTargetSettingsBtn.ConnectClicked(func(bool) {
		qtui.DumperWhitelistWindow.Show()
	})

	qtui.Main.ButtonParametersSettings.ConnectClicked(func(bool){
		qtui.ParametersSettingsWindow.Show()
	})

	return qtui.Main
}