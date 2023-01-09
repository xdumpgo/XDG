package uiqt

import (
	"fmt"
	"github.com/xdumpgo/XDG/injection"
	"github.com/xdumpgo/XDG/modules"
	"github.com/xdumpgo/XDG/qtui"
	"github.com/xdumpgo/XDG/utils"
	"github.com/alecthomas/geoip"
	"github.com/spf13/viper"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"strconv"
	"strings"
)

func NewSingleSiteAnalyzer() *qtui.SingleSiteAnalyzer {
	qtui.SingleSiteWindow = qtui.NewSingleSiteAnalyzer(nil)

	qtui.SingleSiteWindow.SingleDumpStart.ConnectClicked(func(bool) {
		if modules.Analyzer.Injectable == nil {
			qtui.SimpleMB(qtui.SingleSiteWindow, "Please select an injectable first", "Error")
		}
		go modules.Analyzer.DumpSelection()
	})

	qtui.SingleSiteWindow.SingleSiteSelectInjection.ConnectClicked(func(checked bool) {
		utils.Done = make(chan interface{})
		go func() {
			if modules.Analyzer.Injectable = injection.ParseUrlString(strings.ReplaceAll(qtui.SingleSiteWindow.SingleSiteText.Text(), "\n", "")); modules.Analyzer.Injectable == nil {
				go qtui.SimpleMB(qtui.SingleSiteWindow, "Please load an injectable", "Error").Show()
				return
			}

			modules.Analyzer.Injectable.GetChunkSize()

			qtui.SingleSiteWindow.SingleVector.SetText(injection.GetVectorName(modules.Analyzer.Injectable.Vector))
			qtui.SingleSiteWindow.SingleTechnique.SetText(injection.TechniqueString(modules.Analyzer.Injectable.Technique))
			qtui.SingleSiteWindow.SingleDatabaseType.SetText(injection.DBMSString(modules.Analyzer.Injectable.DBType))

			ip, err := modules.Analyzer.Injectable.GetIP()
			if err == nil {
				geo, err := geoip.New()
				if err == nil {
					if country := geo.Lookup(ip); country != nil {
						qtui.SingleSiteWindow.SingleCountry.SetText(country.Long)
					}
				}
			}

			qtui.SingleSiteWindow.SingleDatabaseVersion.SetText("Getting...")
			if version, err := modules.Analyzer.Injectable.GetDBVersion(); err == nil {
				qtui.SingleSiteWindow.SingleDatabaseVersion.SetText(version)
			} else {
				qtui.SingleSiteWindow.SingleDatabaseVersion.SetText("Failed")
			}
			qtui.SingleSiteWindow.SingleDatabaseVersion.Repaint()

			qtui.SingleSiteWindow.SingleCurrentUser.SetText("Getting...")
			if user, err := modules.Analyzer.Injectable.GetDBUser(); err == nil {
				qtui.SingleSiteWindow.SingleCurrentUser.SetText(user)
			} else {
				qtui.SingleSiteWindow.SingleCurrentUser.SetText("Failed")
			}
			qtui.SingleSiteWindow.SingleCurrentUser.Repaint()
		}()
	})

	qtui.SingleSiteWindow.SingleGatherStructure.ConnectClicked(func(bool) {
		go modules.Analyzer.GetStructure()
	})

	qtui.SingleSiteWindow.SingleGatherStructureCancel.ConnectClicked(func(checked bool) {
		select {
		case <-utils.Done:
			return
		default:
			close(utils.Done)
		}
	})

	qtui.SingleSiteWindow.SingleDatabaseStructure.ConnectItemClicked(func(item *widgets.QTreeWidgetItem, column int) {
		if column == 2 {
			fmt.Println("checked", item.CheckState(column), item.Parent().Text(0))
			setVal := false
			setChk := core.Qt__Unchecked
			if item.CheckState(2) == 2 {
				setVal = true
				setChk = core.Qt__Checked
			}
			type change struct {
				Item *widgets.QTreeWidgetItem
				DbIndex int
				TbIndex int
				ColIndex int
				Verdict bool
			}


			if item.Parent().Text(0) == "" { // We are DB
				for index, _ := range modules.Analyzer.Injectable.Structure {
					dbItem := &modules.Analyzer.Injectable.Structure[index]
					if dbItem.Name == item.Text(0) {
						dbItem.Selected = setVal
						fmt.Println(item.Text(0), "=>", setChk, "|", strconv.FormatBool(dbItem.Selected))
						for tbIndex, _ := range modules.Analyzer.Injectable.Structure[index].Tables {
							tbItem := &modules.Analyzer.Injectable.Structure[index].Tables[tbIndex]
							c := item.Child(tbIndex)
							tbItem.Selected = setVal
							c.SetCheckState(2, setChk)
							fmt.Println(c.Text(0), "=>", setChk, "|", strconv.FormatBool(tbItem.Selected))
							for cIndex, _ := range tbItem.Columns {
								cItem := &modules.Analyzer.Injectable.Structure[index].Tables[tbIndex].Columns[cIndex]
								cItem.Selected = setVal
								c.Child(cIndex).SetCheckState(2, setChk)
								fmt.Println(c.Child(cIndex).Text(0), "=>", setChk, "|", strconv.FormatBool(cItem.Selected))
							}
						}
					}
				}
			} else if item.Parent().Parent().Text(0) == "" {
				for index, _ := range modules.Analyzer.Injectable.Structure {
					//if dbItem.Name == item.Parent().Text(0) {
					dbItem := &modules.Analyzer.Injectable.Structure[index]
					atLeastOneTB := false
					for tbIndex, _ := range modules.Analyzer.Injectable.Structure[index].Tables {
						tbItem := &modules.Analyzer.Injectable.Structure[index].Tables[tbIndex]
						if tbItem.Name == item.Text(0) {
							dbItem.Selected = setVal
							item.Parent().SetCheckState(2, setChk)
							item.SetCheckState(2, setChk)
							fmt.Println(item.Text(0), "=>", setChk, "|", strconv.FormatBool(tbItem.Selected))
							for cIndex, _ := range tbItem.Columns {
								cItem := &modules.Analyzer.Injectable.Structure[index].Tables[tbIndex].Columns[cIndex]
								cItem.Selected = setVal
								item.Child(cIndex).SetCheckState(2, setChk)
								fmt.Println(item.Child(cIndex).Text(0), "=>", setChk, "|", strconv.FormatBool(cItem.Selected))
							}
						}
						if tbItem.Selected {
							atLeastOneTB = true
						}
					}
					if dbItem != nil && atLeastOneTB {
						dbItem.Selected = true
						item.Parent().SetCheckState(2, core.Qt__Checked)
					}
					//}
				}
			} else {
				for index, _ := range modules.Analyzer.Injectable.Structure {
					atLeastOneTB := false
					dbItem := &modules.Analyzer.Injectable.Structure[index]
					//if dbItem.Name == item.Parent().Parent().Text(0) {
						for tbIndex, _ := range modules.Analyzer.Injectable.Structure[index].Tables {
							atLeastOneCol := false
							tbItem := &modules.Analyzer.Injectable.Structure[index].Tables[tbIndex]
							//if tbItem.Name == item.Parent().Text(0) {
								for cIndex, _ := range tbItem.Columns {
									cItem := &modules.Analyzer.Injectable.Structure[index].Tables[tbIndex].Columns[cIndex]
									if cItem.Name == item.Text(0) {
										item.Parent().Parent().SetCheckState(2, setChk)
										item.Parent().SetCheckState(2, setChk)
										dbItem.Selected = setVal
										tbItem.Selected = setVal
										cItem.Selected = setVal
									}
									if cItem.Selected {
										atLeastOneCol = true
									}
								}
							//}
							if tbItem != nil && atLeastOneCol {
								tbItem.Selected = true
								atLeastOneTB = true
								item.Parent().SetCheckState(2, core.Qt__Checked)
							}
						}
					//}
					if dbItem != nil && atLeastOneTB {
						dbItem.Selected = true
						item.Parent().Parent().SetCheckState(2, core.Qt__Checked)
					}
				}
			}
		}
	})

	viper.SetDefault("analyzer.skipinfoschema", true)
	qtui.SingleSiteWindow.SingleSkipInformationSchema.SetChecked(viper.GetBool("analyzer.skipinfoschema"))
	qtui.SingleSiteWindow.SingleSkipInformationSchema.ConnectClicked(func(checked bool) {
		viper.Set("analyzer.skipinfoschema", checked)
		viper.WriteConfig()
	})

	qtui.SingleSiteWindow.SingleDumpOutput.ConnectItemDoubleClicked(func(item *widgets.QTableWidgetItem) {
		row := item.Row()
		if win, ok := qtui.DumpLogMap[row]; ok {
			//core.QMetaObject_InvokeMethod(win, "Show", core.Qt__QueuedConnection, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
			//win.InvokeMethod("Show", nil)
			win.ShowNormal()
		}
	})

	viper.SetDefault("analyzer.workers", 35)
	qtui.SingleSiteWindow.SpinBox.SetValue(viper.GetInt("analyzer.workers"))
	qtui.SingleSiteWindow.SpinBox.SetMaximum(80)
	qtui.SingleSiteWindow.SpinBox.SetMinimum(1)
	qtui.SingleSiteWindow.SpinBox.ConnectValueChanged(func(i int) {
		viper.Set("analyzer.worker", i)
		viper.WriteConfig()
	})

	return qtui.SingleSiteWindow
}