package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/xdumpgo/XDG/api/server"
	main2 "github.com/xdumpgo/XDG/cmd/apiserver"
	_ "github.com/go-sql-driver/mysql"
	"net"
	"os"
)

const (
	CONN_HOST = "37.228.132.179"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

var (
	AESKEY = []byte{0x4F, 0x9E, 0xDC, 0xE9, 0x15, 0x14, 0x2E, 0xDA, 0xAE, 0xFB, 0xA5, 0x6D, 0x51, 0x18, 0xCC, 0xA5, 0xF6, 0x95, 0x04, 0x9F, 0x16, 0x53, 0x48, 0x13, 0xC1, 0xE4, 0x83, 0x56, 0x1E, 0xB8, 0x53, 0x2C}
	DatabaseConnection *sql.DB
)

func main() {
	DatabaseConnection = CreateCon()
	defer DatabaseConnection.Close()

	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func CreateCon() *sql.DB {
	db, err := sql.Open("mysql", "qauth:fA7FanTBZk^cHLf8@tcp(localhost:3306)/quartzauth?parseTime=true")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	} else {
		fmt.Println("db is connected")
	}
	//defer db.Close()
	// make sure connection is available
	err = db.Ping()
	fmt.Println(err)
	if err != nil {
		fmt.Println("MySQL db is not connected")
		fmt.Println(err.Error())
	}
	return db
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	message, _ := bufio.NewReader(conn).ReadString('\n')
	//	fmt.Println("[original from client]",message)
	decrypted := server.Decrypt(message, "zLy&CyLg#tvUp4aH6tDH7F%fx=w6xCa%r_A^3hkq468nuuc5=7xba^un&bYLeQ=C-qL_hp#rnNU5a!5neb%_&aygXDL8Jg7?Y6cHpLAtEXM&kbfGTzyQm3Lv")
	var fc *server.FromClient
	fmt.Println("[after decrypt]",decrypted)
	err := json.Unmarshal([]byte(decrypted), &fc)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fc = server.DecryptPayload(fc)
	//fmt.Println(fc)
	str, _ := json.Marshal(fc)
	fmt.Println("[after payload decrypt]",string(str))
	tc := main2.validate(fc)
	if tc.Status == "success" {
		tc.IP = conn.RemoteAddr().String()

		switch fc.PacketType {
		case "heartbeat": // heartbeat
			tc = main2.heartbeat(fc)
		case "authenticate": // Login
			tc = main2.login(fc)
		case "register": // register
			tc = main2.register(fc)
		case "redeem": // renew
			tc = main2.renew(fc)
		case "var":
			tc = main2.nvar(fc)
		case "avar":
			tc = main2.avar(fc)
		default:
			tc = server.GenError("Unknown Error")
		}
	} else {
		fmt.Println("Invalid request")
	}

	fmt.Println("[before encrypt] ip: ", conn.RemoteAddr().String() ,tc)
	resp, err := json.Marshal(server.EncryptPayload(tc, fc.Session, fc.Salt))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//	fmt.Println("[Encrypt payload]", string(resp))
	raw := server.Encrypt(string(resp), "zLy&CyLg#tvUp4aH6tDH7F%fx=w6xCa%r_A^3hkq468nuuc5=7xba^un&bYLeQ=C-qL_hp#rnNU5a!5neb%_&aygXDL8Jg7?Y6cHpLAtEXM&kbfGTzyQm3Lv")
	//	fmt.Println("[sending]",raw)
	// Send a response back to person contacting us.
	conn.Write([]byte(raw))
}
