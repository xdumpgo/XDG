package uiqt

import (
	"fmt"
	"github.com/xdumpgo/XDG/auth"
	"github.com/xdumpgo/XDG/injection"
	"github.com/xdumpgo/XDG/qtui"
	"github.com/xdumpgo/XDG/utils"
	"github.com/pkg/browser"
	"github.com/therecipe/qt/widgets"
	"io/ioutil"
	"os"
	"time"
)

func NewAuthWindow() *qtui.AuthForm {
	qtui.AuthWindow = qtui.NewAuthForm(nil)

	qtui.AuthWindow.LoginBtn.ConnectClicked(func(bool) {
		resp := auth.Login(qtui.AuthWindow.LoginUsername.Text(), qtui.AuthWindow.LoginPassword.Text())

		if resp.Status == "success" {
			msg := widgets.NewQMessageBox(qtui.AuthWindow)
			msg.SetWindowTitle("Authentication Success")
			msg.SetText("Welcome, " + auth.ClientUsername)
			msg.Show()
			msg.ActivateWindow()
			msg.ConnectAccepted(func() {
				if qtui.AuthWindow.LoginRememberMe.IsChecked() {
					ioutil.WriteFile(".xdg", []byte(utils.Encrypt(fmt.Sprintf("%s:%s", auth.ClientUsername, auth.Password), "xNz#'%/2n4SZsB>m")), 0644)
				}
				injection.Init()
				qtui.AuthWindow.Hide()
				qtui.Main.Show()
				go func() {
					for {
						select {
						case <- time.After(30 * time.Second):
							if auth.Expiry.Before(time.Now()) {
								m := qtui.SimpleMB(qtui.Main, "Looks like your license expired, get a new one from https://xdg.quartzinc.dev/", "Uh oh")
								close(utils.Done)
								close(utils.Kill)
								m.ConnectAccepted(func() {
									os.Exit(0)
								})
								m.Show()
								m.ActivateWindow()
								qtui.Main.Hide()
							}
						}
					}
				}()
			})
		} else {
			msg := widgets.NewQMessageBox(qtui.AuthWindow)
			msg.SetWindowTitle("Authentication Failure")
			switch resp.Message {
			case "invalid_credentials":
				msg.SetText("Invalid Credentials")
			case "invalid_hwid":
				msg.SetText("Invalid HWID")
			case "license_expired":
				msg.SetText("Expired License")
			case "update_available":
				msg.SetText("Update Available, press OK to download")
				msg.ConnectAccepted(func() {
					browser.OpenURL(resp.Data)
				})
			}
			msg.Show()
			msg.ActivateWindow()
		}
	})

	qtui.AuthWindow.RegisterBtn.ConnectClicked(func(bool) {
		resp := auth.Register(qtui.AuthWindow.RegisterUsername.Text(), qtui.AuthWindow.RegisterEmail.Text(), qtui.AuthWindow.RegisterPassword.Text(), qtui.AuthWindow.RegisterToken.Text())

		if resp.Status == "success" {
			msg := widgets.NewQMessageBox(qtui.AuthWindow)
			msg.SetWindowTitle("Register Success")
			msg.SetText("Thank you for registering!  Please login")
			msg.Show()
			msg.ActivateWindow()
		} else {
			msg := widgets.NewQMessageBox(qtui.AuthWindow)
			msg.SetWindowTitle("Register Failure")
			switch resp.Message {
			case "update_available":
				msg.SetText("Update Available, press OK to download")
				msg.ConnectAccepted(func() {
					browser.OpenURL(resp.Data)
				})
			case "invalid_token":
				msg.SetText("Invalid token")
			}
			msg.Show()
			msg.ActivateWindow()
		}
	})

	qtui.AuthWindow.RedeemBtn.ConnectClicked(func(bool) {
		resp := auth.Renew(qtui.AuthWindow.RedeemUsername.Text(), qtui.AuthWindow.RedeemPassword.Text(), qtui.AuthWindow.RedeemToken.Text())

		if resp.Status == "success" {
			msg := widgets.NewQMessageBox(qtui.AuthWindow)
			msg.SetWindowTitle("Redeem Success")
			msg.SetText("Successfully redeemed token")
			msg.Show()
			msg.ActivateWindow()
		} else {
			msg := widgets.NewQMessageBox(qtui.AuthWindow)
			msg.SetWindowTitle("Redeem Failure")
			msg.SetText("Invalid Token")
			msg.Show()
			msg.ActivateWindow()
		}
	})

	return qtui.AuthWindow
}