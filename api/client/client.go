package client

import (
	"fmt"
	protocol "github.com/xdumpgo/XDG/apiproto"
	"github.com/xdumpgo/XDG/utils"
)

type messageHandler func(string)

type ChatClient interface {
	Dial(address string) error
	Start()
	Close()
	Send(command interface{}) error
	SetName(name string) error
	SendMessage(message string) error
	Incoming() chan protocol.MessageCommand
	UserList() chan protocol.UserList
	Terminate() chan protocol.Terminate
}

var XDGAPI *TcpClient

var (
	AuthenticationServers = []string {
		"95.217.194.19",
		"95.216.108.152",
		"95.216.108.153",
		"95.216.108.154",
		"95.216.108.155",
		"95.216.108.156",
		"95.216.108.157",
		"95.216.108.158",
		"95.216.108.159",
	}
	GoodServer string
)

func ConnectToAPIServer() error {
	XDGAPI = NewClient()

	for i, ip := range AuthenticationServers {
		utils.LogInfo(fmt.Sprintf("Attempting connection to license server #%d\n", i))
		GoodServer = fmt.Sprintf("%s:7005", ip)
		err := XDGAPI.Dial(GoodServer)
		if err == nil {
			break
		}
	}

	go XDGAPI.Start()

	return nil
}