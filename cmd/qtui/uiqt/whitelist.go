package uiqt

import (
	"bufio"
	"fmt"
	"github.com/xdumpgo/XDG/modules"
	"github.com/xdumpgo/XDG/qtui"
	"github.com/xdumpgo/XDG/utils"
	"github.com/spf13/viper"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"strings"
)

func NewWhitelistWindow() *qtui.DumperWhitelist {
	qtui.DumperWhitelistWindow = qtui.NewDumperWhitelist(nil)

	viper.SetDefault("dumper.whitelist", map[string][]string{
		"primary": {
			"user",
			"mail",
			"login",
		},
		"secondary": {
			"pass",
			"pssw",
			"pwd",
			"hash",
			"md5",
		},
	})

	modules.Dumper.MarshalWhiteList(viper.GetStringMapStringSlice("dumper.whitelist"))
	for _, group := range modules.Dumper.Whitelist {
		qtui.DumperWhitelistWindow.WhitelistGroups.AddItem(group.Name)
	}

	qtui.DumperWhitelistWindow.WhitelistGroups.Item(0).SetSelected(true)
	qtui.DumperWhitelistWindow.WhitelistData.SetPlainText(strings.Join(modules.Dumper.Whitelist[0].List, "\n"))
	qtui.DumperWhitelistWindow.WhitelistGroups.ConnectSelectionChanged(func(selected *core.QItemSelection, deselected *core.QItemSelection) {
		qtui.DumperWhitelistWindow.WhitelistData.SetPlainText(strings.Join(modules.Dumper.Whitelist[qtui.DumperWhitelistWindow.WhitelistGroups.CurrentRow()].List, "\n"))
	})

	qtui.DumperWhitelistWindow.WhitelistOk.ConnectClicked(func(checked bool) {
		viper.Set("dumper.whitelist", modules.Dumper.UnmarshalWhitelist())
		viper.WriteConfig()
		qtui.DumperWhitelistWindow.Close()
	})

	qtui.DumperWhitelistWindow.WhitelistCancel.ConnectClicked(func(checked bool) {
		modules.Dumper.MarshalWhiteList(viper.GetStringMapStringSlice("dumper.whitelist"))
		qtui.DumperWhitelistWindow.Close()
	})

	qtui.DumperWhitelistWindow.WhitelistAddGroup.ConnectClicked(func(checked bool) {
		add := qtui.NewWhitelistAdd(nil)
		add.CreateBtn.ConnectClicked(func(checked bool) {
			for _, group := range modules.Dumper.Whitelist {
				if group.Name == add.GroupName.Text() {
					go qtui.SimpleMB(qtui.DumperWhitelistWindow, "You already have a group with that name", "Error").Show()
					return
				}
			}
			modules.Dumper.Whitelist = append(modules.Dumper.Whitelist, modules.DumperWhitelistGroup{
				Name: add.GroupName.Text(),
				List: make([]string, 0),
			})
			qtui.DumperWhitelistWindow.WhitelistGroups.AddItem(add.GroupName.Text())

			add.Hide()
			add.Close()
		})
		add.CancelBtn.ConnectClicked(func(checked bool) {
			add.Hide()
			add.Close()
		})
		add.Show()
	})

	viper.SetDefault("dumper.blacklist", []string{"id", "last", "date", "time"})
	qtui.DumperWhitelistWindow.Blacklist.SetPlainText(strings.Join(viper.GetStringSlice("dumper.blacklist"), "\n"))
	qtui.DumperWhitelistWindow.Blacklist.ConnectFocusOutEvent(func(event *gui.QFocusEvent) {
		viper.Set("dumper.blacklist", strings.Split(qtui.DumperWhitelistWindow.Blacklist.ToPlainText(), "\n"))
		viper.WriteConfig()
	})

	qtui.DumperWhitelistWindow.WhitelistRemoveGroup.ConnectClicked(func(checked bool) {
		if len(modules.Dumper.Whitelist) == 1 {
			go qtui.SimpleMB(qtui.DumperWhitelistWindow, "You must have one Whitelist Group", "Error").Show()
			return
		}
		index := qtui.DumperWhitelistWindow.WhitelistGroups.CurrentRow()
		modules.Dumper.Whitelist = remove(modules.Dumper.Whitelist, index)
		qtui.DumperWhitelistWindow.WhitelistGroups.Item(index).DestroyQListWidgetItem()
	})

	qtui.DumperWhitelistWindow.WhitelistData.ConnectTextChanged(func() {
		scanner := bufio.NewScanner(strings.NewReader(qtui.DumperWhitelistWindow.WhitelistData.ToPlainText()))
		index := qtui.DumperWhitelistWindow.WhitelistGroups.CurrentRow()
		modules.Dumper.Whitelist[index].List = make([]string, 0)
		for scanner.Scan() {
			if !utils.ArrContains(modules.Dumper.Whitelist[index].List, scanner.Text()) {
				modules.Dumper.Whitelist[index].List = append(modules.Dumper.Whitelist[index].List, scanner.Text())
			}
		}
	})

	return qtui.DumperWhitelistWindow
}

func remove(slice []modules.DumperWhitelistGroup, s int) []modules.DumperWhitelistGroup {
	fmt.Println(slice)
	_new := append(slice[:s], slice[s+1:]...)
	fmt.Printf("%#v", _new)
	return _new
}