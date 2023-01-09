package auth

import
(
	"time"
)

type ToServer struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email string `json:"email"`
	HWID string `json:"hwid"`
	PacketType string `json:"packettype"`
	Hash string `json:"hash"`
	Data string `json:"data"`
	Key string `json:"key"`
	VKey string `json:"vkey"`
	Version string `json:"version"`
	Session string `json:"session_id"`
	Salt string `json:"salt"`
}

type FromServer struct {
	Status string `json:"status"`
	Message string `json:"message"`
	Data string `json:"data"`
	ArrData map[string]string `json:"adata"`
	Expiry time.Time `json:"expiry"`
	Level int `json:"level"`
	IP string `json:"ip"`
}
