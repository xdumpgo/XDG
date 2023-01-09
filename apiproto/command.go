package apiproto

import (
	"errors"
	"time"
)

var (
	UnknownCommand = errors.New("unknown command")
)

const (
	HEARTBEAT = 0x00
	SEND = 0x01
	MESSAGE = 0x02
	USERLIST = 0x03
	LOGIN = 0x04
	REGISTER = 0x05
	REDEEM = 0x06
	VAR = 0x07
	RESPONSE = 0x08
	PROXYLIST = 0x09
	NEWS = 0x10
	STATUS = 0x11
	NAME = 0x12

	MELT = 0x98
	TERMINATE = 0x99
)

type SessionData struct {
	SessionID string `json:"session_id"`
	SessionSalt string `json:"session_salt"`
}

// SendCommand is used for sending new message from client
type SendCommand struct {
	SessionData
	Message string `json:"message"`
}

// MessageCommand is used for notifying new messages
type MessageCommand struct {
	SessionData
	Name    string `json:"name"`
	Message string `json:"message"`
}

type Proxies struct {
	SessionData
	List []string `json:"list"`
}

// UserListCommand
type UserList struct {
	SessionData
	Users []string `json:"users"`
}

type Terminate struct {
	SessionData
	Reason string `json:"reason"`
}

type StatsUpdate struct {
	SessionData
	CurrentModule string `json:"current_module"`
	Runtime time.Duration `json:"runtime"`
	Index int `json:"index"`
	End int `json:"end"`
	Threads int `json:"threads"`
	Workers int `json:"workers"`
	Urls int `json:"urls"`
	Injectables int `json:"injectables"`
	Rows int `json:"rows"`
}

type AuthProto struct {
	SessionData
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	HWID string `json:"hwid,omitempty"`
	Version string `json:"version,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type AuthResponse struct {
	Status string `json:"status"`
	Message string `json:"message,omitempty"`
	Data string `json:"data,omitempty"`
	ArrData map[string]string `json:"arr_data,omitempty"`
	Expiry time.Time `json:"expiry,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type LoginCommand struct {
	AuthProto 
}

type RegisterCommand struct {
	AuthProto 
	Email string `json:"email"`
	Token string `json:"token"`
}

type RedeemCommand struct {
	AuthProto 
	Token string `json:"token"`
}

type VarCommand struct {
	AuthProto 
	Name string `json:"name"`
}

type VarDumpCommand struct {
	AuthProto
	Names []string `json:"names"`
}

type Heartbeat struct {
	AuthProto 
}

type News struct {
	Message string
	Time time.Time
}

type NewsCommand struct {
	AuthProto
	News []News
}

type Melt struct {

}

type NameCommand struct {
	Name string `json:"name"`
}