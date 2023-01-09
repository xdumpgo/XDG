package auth

import (
	"fmt"
	"github.com/xdumpgo/XDG/api/client"
	"github.com/xdumpgo/XDG/apiproto"
	"github.com/xdumpgo/XDG/utils"
	"time"
)

var
(
	Version = "1.5.5h3"
	Expiry time.Time
	ClientUsername = ""
	Password = ""
)

func Login(username string, password string) apiproto.AuthResponse {
	utils.LogInfo("Authenticating..")
	if err := client.XDGAPI.Login(username, password, Version); err != nil {
		fmt.Println(err.Error())
		return apiproto.AuthResponse{}
	}

	resp := <- client.XDGAPI.Auth()

	ClientUsername = username
	Password = password
	Expiry = resp.Expiry

	return resp
}

func Register(username string, email string, password string, token string) apiproto.AuthResponse {
	err := client.XDGAPI.Register(username, password, email, token, Version)
	if err != nil {
		return apiproto.AuthResponse{}
	}

	return <- client.XDGAPI.Auth()
}

func Renew(username string, password string, token string) apiproto.AuthResponse {
	err := client.XDGAPI.Redeem(username, password, token, Version)
	if err != nil {
		return apiproto.AuthResponse{}
	}

	return <- client.XDGAPI.Auth()
}

func Var(name string) string {
	err := client.XDGAPI.Var(ClientUsername, Password, name, Version)
	if err != nil {
		return ""
	}
	fmt.Println("Waiting on response")
	if resp := <- client.XDGAPI.Auth(); resp.Status == "success" {
		return resp.Data
	}
	return ""
}

func AVar() map[string]string {
	err := client.XDGAPI.Var(ClientUsername, Password, "all", Version)
	if err != nil {
		return nil
	}


	if resp := <- client.XDGAPI.Auth(); resp.Status == "success" {
		return resp.ArrData
	}
	return nil
}

func Heartbeat() bool {
	err := client.XDGAPI.Heartbeat()
	if err != nil {
		return false
	}

	resp := <- client.XDGAPI.Auth()

	if resp.Status == "failure" {
		utils.LogError("Somethings fishy...")
		return false
	}
	return true
}