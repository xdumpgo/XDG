package server

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/xdumpgo/XDG/apiproto"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

func heartbeat() apiproto.AuthResponse {
	return apiproto.AuthResponse{
		Status:    "success",
		Message:   "alive",
		Timestamp: time.Now(),
	}
}

func login(command apiproto.LoginCommand) apiproto.AuthResponse {
	program, err := GetProgram("iMH5RjESHkMJr8EwbcNjYtcsz9fIXEGlp0fXDc5u")
	if err != nil {
		return GenError("database_error")
	}
	if len(command.Username) == 0 || len(command.Password) == 0 {

		return GenError("invalid_credentials")
	}

	user, err := getUser(command.Username, program.Id)
	if err != nil {
		return GenError("invalid_user")
	}

	if res := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(command.Password)); res != nil {
		return GenError("invalid_credentials")
	}

	if len(user.HWID) == 0 {
		user.SetHWID(command.HWID)
	} else if command.HWID != user.HWID {
		return GenError("invalid_hwid")
	}

	if user.Expires.Before(time.Now()) {
		return GenError("license_expired")
	}

	if command.Version != program.Version  {
		tc := GenError("update_available")
		tc.Data = program.Url.String
		return tc
	}

	return apiproto.AuthResponse{
		Status:  "success",
		Message: "authenticated",
		Expiry:  user.Expires,
		Timestamp: time.Now(),
	}
}

func register(command apiproto.RegisterCommand) apiproto.AuthResponse {
	program, err := GetProgram("iMH5RjESHkMJr8EwbcNjYtcsz9fIXEGlp0fXDc5u")
	if err != nil {
		return GenError(err.Error())
	}
	token, err := getToken(command.Token, program.Id)
	if err != nil {
		return GenError(err.Error())
	}

	stmt, err := DatabaseConnection.Prepare("SELECT COUNT(*) FROM program_users WHERE username=? AND program_id = ?")
	if err != nil {
		return GenError("database_error")
	}
	defer stmt.Close()
	row := stmt.QueryRow(command.Username, program.Id)
	var c int
	err = row.Scan(&c)
	if err != nil {
		return GenError("database_error")
	}

	if c != 0 {
		return GenError("username_in_use")
	}

	stmt, err = DatabaseConnection.Prepare("INSERT INTO program_users (username, password, email, hwid, level, expires, program_id, last_reset) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return GenError("database_error")
	}

	_, err = stmt.Exec(command.Username, hashAndSalt([]byte(command.Password)), command.Email, command.HWID, token.Level, time.Now(), program.Id, time.Now().AddDate(0,0,-7))
	if err != nil {
		return GenError("database_error")
	}

	row = DatabaseConnection.QueryRow("SELECT id FROM program_users WHERE program_id = 57 AND username = ?", command.Username)
	var id int
	row.Scan(&id)

	token.Use(&UserAccount{Id: id, Username: command.Username})

	return apiproto.AuthResponse{
		Status:  "success",
		Message: "registered",
		Expiry: time.Now().AddDate(0,0, token.Days),
		Timestamp: time.Now(),
	}
}

func redeem(command apiproto.RedeemCommand) apiproto.AuthResponse {
	program, err := GetProgram("iMH5RjESHkMJr8EwbcNjYtcsz9fIXEGlp0fXDc5u")
	if err != nil {
		return GenError(err.Error())
	}
	l := apiproto.LoginCommand{}
	l.Username = command.Username
	l.Password = command.Password
	l.HWID = command.HWID
	l.Version = command.Version
	tc := login(l)

	token, err := getToken(command.Token, program.Id)
	if err != nil {
		return GenError(err.Error())
	}

	user, err := getUser(command.Username, program.Id)

	if tc.Status == "success" {
		t, err := token.Use(user)
		if err != nil {
			return GenError("invalid_token")
		}
		return apiproto.AuthResponse{
			Status:    "success",
			Message:   "token_redeemed",
			Expiry:    *t,
			Timestamp: time.Now(),
		}
	}
	return tc
}

func proxies(command apiproto.Proxies) apiproto.Proxies {
	f, err := os.Open("proxies.txt")
	if err != nil {
		return apiproto.Proxies{}
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		command.List = append(command.List, scanner.Text())
	}
	return command
}

func nvar(command apiproto.VarCommand) apiproto.AuthResponse {
	if command.Name == "all" {
		fmt.Println("getting all")
		return avar(command)
	}

	program, err := GetProgram("iMH5RjESHkMJr8EwbcNjYtcsz9fIXEGlp0fXDc5u")
	if err != nil {
		return GenError(err.Error())
	}

	l := apiproto.LoginCommand{}
	l.Username = command.Username
	l.Password = command.Password
	l.HWID = command.HWID
	l.Version = command.Version
	tc := login(l)

	if tc.Status != "success" {
		return tc
	}

	//if fc.VKey != program.VarKey {
	//	return GenError("invalid_variable_key")
	//}

	stmt, err := DatabaseConnection.Prepare("SELECT `value` FROM variables WHERE program_id = ? AND Name = ? LIMIT 1");
	if err != nil {
		return GenError("database error")
	}
	defer stmt.Close()
	tc = apiproto.AuthResponse{Status: "success"}

	row := stmt.QueryRow(program.Id, command.Name)
	switch err := row.Scan(&tc.Data); err {
	case sql.ErrNoRows:
		return GenError("invalid_variable")
	case nil:
		return tc
	default:
		panic(err)
	}
}

func avar(command apiproto.VarCommand) apiproto.AuthResponse {
	program, err := GetProgram("iMH5RjESHkMJr8EwbcNjYtcsz9fIXEGlp0fXDc5u")
	if err != nil {
		return GenError(err.Error())
	}

	l := apiproto.LoginCommand{}
	l.Username = command.Username
	l.Password = command.Password
	l.HWID = command.HWID
	l.Version = command.Version
	tc := login(l)

	if tc.Status != "success" {
		return tc
	}

	stmt, err := DatabaseConnection.Prepare("SELECT `Name`,`value` FROM variables WHERE program_id = ?");
	if err != nil {
		return GenError("database error")
	}
	defer stmt.Close()

	rows, err := stmt.Query(program.Id)
	if err != nil {
		return GenError("database_error")
	}
	tc.ArrData = make(map[string]string)

	for rows.Next() {
		var name string
		var val string
		err = rows.Scan(&name, &val)
		if err != nil {
			panic(err)
		}
		tc.ArrData[name] = val
	}
	return tc
}
