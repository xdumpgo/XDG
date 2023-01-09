package server

import (
	"crypto/hmac"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/xdumpgo/XDG/apiproto"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"time"
)

var (
	DatabaseConnection *sql.DB
)

func GenError(message string) apiproto.AuthResponse {
	return apiproto.AuthResponse{
		Status:  "error",
		Message: message,
		Expiry:  time.Time{},
	}
}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, 14)
	if err != nil {
		log.Println(err)
	}    // GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func getToken(token string, program int) (*Token, error) {
	t := &Token{Token: token, ProgramID: program}
	stmt, err := DatabaseConnection.Prepare("SELECT level, days, used FROM tokens WHERE token = ? AND program_id = ?")
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("database_error")
	}
	defer stmt.Close()
	row := stmt.QueryRow(token, program)
	switch err := row.Scan(&t.Level, &t.Days, &t.Used); err {
	case sql.ErrNoRows:
		fmt.Println("not exist")
		return nil, errors.New("invalid_token")
	case nil:
		if t.Used == 1 {
			return nil, errors.New("invalid_token")
		}
		return t, nil
	default:
		panic(err)
	}
}


func GetProgram(key string) (*Program, error) {
	p := &Program{Key: key}
	stmt, err := DatabaseConnection.Prepare("SELECT Id, Name, hash, varkey, enabled, devmode, url, version FROM programs WHERE `key` = ?")
	if err != nil {
		fmt.Println("stmt error", err.Error())
		return nil, errors.New("database_error")
	}
	defer stmt.Close()
	row := stmt.QueryRow(p.Key)
	switch err := row.Scan(&p.Id, &p.Name, &p.Hash, &p.VarKey, &p.Enabled, &p.Devmode, &p.Url, &p.Version); err {
	case sql.ErrNoRows:
		return nil, errors.New("invalid_program")
	case nil:
		return p, nil
	default:
		panic(err)
	}
}

func getProgramByOwnerID(id int) (*Program, error) {
	p := &Program{OwnerID: id}
	stmt, err := DatabaseConnection.Prepare("SELECT Id, Name, hash, key, varkey, enabled, devmode, url, version FROM programs WHERE `owner_id` = ?")
	if err != nil {
		fmt.Println("stmt error", err.Error())
		return nil, errors.New("database_error")
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	switch err := row.Scan(&p.Id, &p.Name, &p.Hash, &p.Key, &p.VarKey, &p.Enabled, &p.Devmode, &p.Url, &p.Version); err {
	case sql.ErrNoRows:
		return nil, errors.New("invalid_program")
	case nil:
		return p, nil
	default:
		panic(err)
	}
}

func getUser(username string, program int) (*UserAccount, error) {
	ua := &UserAccount{Username: username, ProgramID: program}
	stmt, err := DatabaseConnection.Prepare("SELECT Id, password, email, hwid, level, expires FROM program_users WHERE username = ? AND program_id = ?")
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("database_error")
	}
	defer stmt.Close()
	row := stmt.QueryRow(username, program)
	switch err := row.Scan(&ua.Id, &ua.Password, &ua.Email, &ua.HWID, &ua.Level, &ua.Expires); err {
	case sql.ErrNoRows:
		return nil, errors.New("invalid_user")
	case nil:
		return ua, nil
	default:
		panic(err)
	}
}

func GetFrontendLinkUser(id string) int {
	stmt, err := DatabaseConnection.Prepare("SELECT user_id FROM discord_link WHERE discord_id = ?")
	if err != nil {
		return -1
	}
	var user_id int

	row := stmt.QueryRow(id)
	switch err := row.Scan(&user_id); err {
	case sql.ErrNoRows:
		return 0
	case nil:
		return user_id
	default:
		panic(err)
	}
}

func GetProgramsByUser(id int) []*Program {
	stmt, err := DatabaseConnection.Prepare("SELECT program_id, access FROM program_admins WHERE user_id = ?")
	if err != nil {
		fmt.Println("stmt error", err.Error())
		return nil
	}

	rows, err := stmt.Query(id)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var ors string
	ids := make(map[int]int)

	for rows.Next() {
		var ii int
		var access int
		err = rows.Scan(&ii, &access)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		ids[ii] = access
		ors += fmt.Sprintf(" OR id = %d", ii)
	}

	stmt, err = DatabaseConnection.Prepare("SELECT id, name, hash, `key`, varkey, enabled, devmode, url, version FROM programs WHERE `owner_id` = ?" + ors)
	if err != nil {
		fmt.Println("stmt error", err.Error())
		return nil
	}

	var arr []*Program
	rows, err = stmt.Query(id)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	for rows.Next() {
		p := &Program{}
		err = rows.Scan(&p.Id, &p.Name, &p.Hash, &p.Key, &p.VarKey, &p.Enabled, &p.Devmode, &p.Url, &p.Version)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		acc := ids[p.Id]
		if acc != 0 {
			p.Access = acc
		} else {
			p.Access = 3
		}
		arr = append(arr, p)
	}

	return arr
}
/*
func EncryptPayload(tc *ToClient, pass string, salt string) *ToClient {
	tc.Status = EncryptWithSalt(tc.Status, pass, salt)
	tc.Message = EncryptWithSalt(tc.Message, pass, salt)
	tc.Data = EncryptWithSalt(tc.Data, pass, salt)
	return tc
}

func DecryptPayload(fc *FromClient) *FromClient {
	fc.Username = DecryptWithSalt(fc.Username, fc.Session, fc.Salt)
	fc.Password = DecryptWithSalt(fc.Password, fc.Session, fc.Salt)
	fc.Email = DecryptWithSalt(fc.Email, fc.Session, fc.Salt)
	fc.HWID = DecryptWithSalt(fc.HWID, fc.Session, fc.Salt)
	fc.Hash = DecryptWithSalt(fc.Hash, fc.Session, fc.Salt)
	fc.Data = DecryptWithSalt(fc.Data, fc.Session, fc.Salt)
	fc.Key = DecryptWithSalt(fc.Key, fc.Session, fc.Salt)
	fc.Version = DecryptWithSalt(fc.Version, fc.Session, fc.Salt)
	fc.PacketType = DecryptWithSalt(fc.PacketType, fc.Session, fc.Salt)
	return fc
}*/

type Field struct {
	name string
	value string
}


func ArrayContains(arr []string, str string) bool {
	for _, k := range arr {
		if k == str {
			return true
		}
	}
	return false
}

func ArrayReverse(a []string) []string {
	b := a
	for i := len(b)/2-1; i >= 0; i-- {
		opp := len(b)-1-i
		b[i], b[opp] = b[opp], b[i]
	}
	return b
}

func ComputeHmac256(message []byte, secret []byte) string {
	hash := hmac.New(sha512.New, secret)
	hash.Write(message)

	// to lowercase hexits
	return hex.EncodeToString(hash.Sum(nil))

	//return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func delaySecond(n time.Duration) {
	time.Sleep(n * time.Second)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrc(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
