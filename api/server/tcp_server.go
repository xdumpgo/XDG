package server

import (
	"errors"
	"fmt"
	protocol "github.com/xdumpgo/XDG/apiproto"
	"github.com/xdumpgo/XDG/utils"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Client struct {
	Id     string
	conn   net.Conn
	Name   string
	Writer *protocol.CommandWriter
	Status chan protocol.StatsUpdate
}

type TcpServer struct {
	listener net.Listener
	clients  []*Client
	mutex    *sync.Mutex
}

type API interface {
	Listen(address string) error
	Broadcast(command interface{}) error
	Start()
	Close()
	GetClients() []*Client
}

var APIServer API

var (
	UnknownClient = errors.New("unknown client")
)

func NewServer() *TcpServer {
	return &TcpServer{
		mutex: &sync.Mutex{},
	}
}

func (s *TcpServer) GetClientNameArray() []string {
	cl := []string{}
	for _, k := range s.clients {
		cl = append(cl, k.Name)
	}
	return cl
}

func (s *TcpServer) Listen(address string) error {
	l, err := net.Listen("tcp", address)

	if err == nil {
		s.listener = l
	}

	log.Printf("Listening on %v", address)

	return err
}

func (s *TcpServer) Close() {
	s.listener.Close()
}

func (s *TcpServer) Start() {
	for {
		// XXX: need a way to break the loop
		conn, err := s.listener.Accept()

		if err != nil {
			log.Print(err)
		} else {
			// handle connection
			client := s.accept(conn)
			go s.serve(client)
		}
	}
}

func (s *TcpServer) Broadcast(command interface{}) error {
	for _, client := range s.clients {
		// TODO: handle error here?
		client.Writer.Write(command)
	}

	return nil
}

func (s *TcpServer) GetClients() []*Client {
	return s.clients
}

func (s *TcpServer) Send(name string, command interface{}) error {
	for _, client := range s.clients {
		if client.Name == name {
			return client.Writer.Write(command)
		}
	}

	return UnknownClient
}

func (s *TcpServer) accept(conn net.Conn) *Client {
	log.Printf("Accepting connection from %v, total clients: %v", conn.RemoteAddr().String(), len(s.clients)+1)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	client := &Client{
		Id:     utils.StringOfLength(5),
		conn:   conn,
		Writer: protocol.NewCommandWriter(conn),
		Status: make(chan protocol.StatsUpdate),
	}

	s.clients = append(s.clients, client)

	return client
}

func (s *TcpServer) remove(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// remove the connections from clients array
	for i, check := range s.clients {
		if check == client {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
		}
	}

	log.Printf("Closing connection from %v", client.conn.RemoteAddr().String())
	client.conn.Close()
}

func (s *TcpServer) serve(client *Client) {
	cmdReader := protocol.NewCommandReader(client.conn)

	defer func() {
		s.remove(client)
		if len(client.Name) > 0 {
			s.Broadcast(protocol.MessageCommand{
				Name:    client.Name,
				Message: "has left the room.",
			})
			s.Broadcast(protocol.UserList{Users: s.GetClientNameArray()})
		}
	}()

	for {
		cmd, err := cmdReader.Read()

		if err != nil && err != io.EOF {
			if err == protocol.UnknownCommand {
				return
			}
			log.Printf("Read error: %v", err)
		}

		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.Heartbeat:
				client.Writer.Write(heartbeat())
			case protocol.SendCommand:
				go s.Broadcast(protocol.MessageCommand{
					Message: v.Message,
					Name:    client.Name,
				})
			case protocol.LoginCommand:
				ret := login(v)
				client.Writer.Write(ret)
				if ret.Status == "success" {
					utils.LogInfo(fmt.Sprintf("[%s] Authenticated user %s", client.conn.RemoteAddr().String(), v.Username))
					client.Name = v.Username
					go s.Broadcast(protocol.UserList{
						Users: s.GetClientNameArray(),
					})
					go s.Broadcast(protocol.MessageCommand{
						Name:    client.Name,
						Message: "has joined the room.",
					})
					t, _ := time.Parse(time.Stamp, "August 4 12:30:00")
					client.Writer.Write(protocol.NewsCommand{
						AuthProto: protocol.AuthProto{},
						News: []protocol.News{
							{
								Message: "I'll fill this up at some point",
								Time:    t,
							},
						},
					})
				} else {
					utils.LogError(fmt.Sprintf("[%s] User %s failed authentication with reason [%s]", client.conn.RemoteAddr().String(), v.Username, ret.Message))
				}
			case protocol.RegisterCommand:
				res := register(v)
				client.Writer.Write(res)
				if res.Status == "success" {
					utils.LogInfo(fmt.Sprintf("[%s] Registered new user %s", client.conn.RemoteAddr().String(), v.Username))
				} else {
					utils.LogError(fmt.Sprintf("[%s] User %s failed registration with reason [%s]", client.conn.RemoteAddr().String(), v.Username, res.Message))
				}
			case protocol.RedeemCommand:
				res := redeem(v)
				client.Writer.Write(res)
				if res.Status == "success" {
					utils.LogInfo(fmt.Sprintf("[%s] User %s redeemed token [%s]", client.conn.RemoteAddr().String(), v.Username, v.Token))
				} else {
					utils.LogError(fmt.Sprintf("[%s] User %s failed to redeem token [%s] with reason [%s]", client.conn.RemoteAddr().String(), v.Username, v.Token, res.Message))
				}
			case protocol.VarCommand:
				err = client.Writer.Write(nvar(v))
			case protocol.StatsUpdate:
				client.Status <- v
			case protocol.Proxies:
				client.Writer.Write(proxies(v))
			case protocol.NameCommand:
				client.Name = v.Name
				go s.Broadcast(protocol.UserList{
					Users: s.GetClientNameArray(),
				})
				go s.Broadcast(protocol.MessageCommand{
					Name:    client.Name,
					Message: "has joined the room.",
				})
			}
		}

		if err == io.EOF {
			break
		}
	}
}