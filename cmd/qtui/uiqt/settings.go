package uiqt

import (
	"bufio"
	"fmt"
	"github.com/xdumpgo/XDG/manager"
	"github.com/xdumpgo/XDG/qtui"
	"github.com/spf13/viper"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"os"
	"strings"
	"time"
)

func NewSettingsWindow() *qtui.Settings {
	qtui.SettingsWindow = qtui.NewSettings(nil)

	qtui.SettingsWindow.ProxyTypeComboBox.ConnectCurrentIndexChanged(func(index int) {
		if index > 0 {
			qtui.SettingsWindow.LoadProxiesButton.SetEnabled(true)
		} else {
			qtui.SettingsWindow.LoadProxiesButton.SetEnabled(false)
		}
	})

	viper.SetDefault("core.proxyapi", false)
	qtui.SettingsWindow.LoadProxiesFromAPI.SetChecked(viper.GetBool("core.proxyapi"))
	qtui.SettingsWindow.LoadProxiesFromAPI.ConnectClicked(func(checked bool) {
		viper.Set("core.proxyapi", checked)
		viper.WriteConfig()
	})

	viper.SetDefault("core.proxylink", "")
	qtui.SettingsWindow.LineEdit.SetText(viper.GetString("core.proxylink"))
	qtui.SettingsWindow.LineEdit.ConnectFocusOutEvent(func(event *gui.QFocusEvent) {
		viper.Set("core.proxylink", qtui.SettingsWindow.LineEdit.Text())
		viper.WriteConfig()
	})

	qtui.SettingsWindow.LoadProxiesButton.ConnectClicked(func(bool) {
		dir, _ := os.Getwd()
		ofd := widgets.NewQFileDialog2(qtui.SettingsWindow, "Load proxies", dir, "Text Files (*.txt)")
		ofd.Show()
		ofd.ConnectFileSelected(func(file string) {
			manager.PManager.LoadFile(file, qtui.SettingsWindow.ProxyTypeComboBox.CurrentIndex()-1)
			for _, proxy := range manager.PManager.Proxies {
				if len(proxy.Auth) > 0 {
					qtui.SettingsWindow.ProxiesTextbox.AppendPlainText(fmt.Sprintf("%s:%s", proxy.Address, proxy.Auth))
				} else {
					qtui.SettingsWindow.ProxiesTextbox.AppendPlainText(fmt.Sprintf("%s", proxy.Address))
				}
			}
		})
	})

	qtui.SettingsWindow.ProxiesTextbox.ConnectTextChanged(func() {
		scanner := bufio.NewScanner(strings.NewReader(qtui.SettingsWindow.ProxiesTextbox.ToPlainText()))
		manager.PManager.Proxies = []*manager.Proxy{}
		manager.PManager.LoadScanner(scanner, qtui.SettingsWindow.ProxyTypeComboBox.CurrentIndex()-1)
	})

	qtui.SettingsWindow.ThreadsSpinbox.SetValue(viper.GetInt("core.Threads"))
	qtui.SettingsWindow.ThreadsSpinbox.ConnectValueChanged(func(i int) {
		viper.Set("core.Threads", i)
		viper.WriteConfig()
	})

	qtui.SettingsWindow.TimeoutsSpinbox.SetValue(viper.GetInt("core.Timeouts"))
	qtui.SettingsWindow.TimeoutsSpinbox.ConnectValueChanged(func(i int) {
		viper.Set("core.Timeouts", i)
		viper.WriteConfig()
	})

	qtui.SettingsWindow.AutoThreadsCheckBox.Hide()

	qtui.SettingsWindow.BatchModeCheckBox.SetChecked(viper.GetBool("core.BatchMode"))
	qtui.SettingsWindow.BatchModeCheckBox.ConnectClicked(func(checked bool) {
		viper.Set("core.BatchMode", checked)
		viper.WriteConfig()
	})

	qtui.SettingsWindow.LoadCustomerProxies.ConnectClicked(func(checked bool) {
		loading := qtui.NewLoadingWindow(qtui.SettingsWindow)
		loading.SetWindowTitle("Loading Proxies...")
		loading.Label.SetText("Loading 540k proxies, please wait...")
		loading.Show()
		go func() {
			loading.Label.SetText(fmt.Sprintf("Loaded %d Proxies, enjoy!", manager.PManager.LoadCustomerProxies(loading)))
			time.Sleep(5)
			loading.Hide()
			loading.Close()
		}()
	})

	return qtui.SettingsWindow
}