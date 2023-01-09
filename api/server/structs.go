package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

type Program struct {
	Id int
	Name string
	Hash sql.NullString
	Key string
	VarKey string
	Enabled int
	Devmode int
	Url sql.NullString
	Version string
	OwnerID int
	Access int
}

func (p *Program) GetUsers() []*UserAccount {
	stmt, err := DatabaseConnection.Prepare("SELECT id,username,email,hwid,expires FROM program_users WHERE program_id = ?")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	rows, err := stmt.Query(p.Id)
	if err != nil {
		fmt.Println("stmt err")
	}

	var arr []*UserAccount
	for rows.Next() {
		ua := &UserAccount{ProgramID:p.Id}
		err = rows.Scan(&ua.Id, &ua.Username, &ua.Email, &ua.HWID, &ua.Expires)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		arr = append(arr, ua)
	}
	return arr
}

func (p *Program) GetUserByName(name string) (*UserAccount, error) {
	stmt, err := DatabaseConnection.Prepare("SELECT id, email, hwid, expires, last_reset FROM program_users WHERE program_id = ? AND username = ?")
	if err != nil {
		return nil, err
	}

	user := &UserAccount{
		Username:  name,
		ProgramID: p.Id,
	}
	if row := stmt.QueryRow(p.Id, name); row != nil {
		if err = row.Scan(&user.Id, &user.Email, &user.HWID, &user.Expires, &user.LastReset); err == nil {
			return user, nil
		}
	}
	return nil, errors.New("failed to find user")
}

func (p *Program) GetUserById(id int) (*UserAccount, error) {
	stmt, err := DatabaseConnection.Prepare("SELECT username, email, hwid, expires, last_reset FROM program_users WHERE program_id = ? AND id = ?")
	if err != nil {
		return nil, err
	}

	user := &UserAccount{
		Id:  id,
		ProgramID: p.Id,
	}
	if row := stmt.QueryRow(p.Id, id); row != nil {
		if err = row.Scan(&user.Username, &user.Email, &user.HWID, &user.Expires, &user.LastReset); err == nil {
			return user, nil
		}
	}
	return nil, errors.New("failed to find user")
}

func (p *Program) GetUserByDiscord(discordid string) (*UserAccount, error) {
	stmt, err := DatabaseConnection.Prepare("SELECT user_id FROM client_discords WHERE discord_id = ?")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	did, _ := strconv.Atoi(discordid)

	r := stmt.QueryRow(did)
	var userId int
	err = r.Scan(&userId)
	if err != nil {
		fmt.Println("1", err.Error())
		return nil, err
	}

	stmt, err = DatabaseConnection.Prepare("SELECT username, email, hwid, expires, last_reset FROM program_users WHERE program_id = ? AND `id` = ?")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	r = stmt.QueryRow(p.Id, userId)
	if r != nil {
		user := &UserAccount{
			Id:  userId,
			ProgramID: p.Id,
		}
		if err = r.Scan(&user.Username, &user.Email, &user.HWID, &user.Expires, &user.LastReset); err == nil {
			return user, nil
		} else {
			fmt.Println("2", err.Error())
		}
	}
	return nil, errors.New("failed to get user, make sure user is linked")
}

func (p *Program) GetToken(token string) (*Token, error) {
	t := &Token{Token: token, ProgramID: p.Id}
	stmt, err := DatabaseConnection.Prepare("SELECT level, days, used FROM tokens WHERE token = ? AND program_id = ?")
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("database_error")
	}
	defer stmt.Close()
	row := stmt.QueryRow(token, p.Id)
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

type UserAccount struct {
	Id int
	Username string
	Password string
	Email string
	HWID string
	Level int
	Expires time.Time
	LastReset time.Time
	ProgramID int
}

func (ua *UserAccount) SetHWID(hwid string) {
	stmt, err := DatabaseConnection.Prepare("UPDATE program_users SET hwid = ?, last_reset = ? WHERE username = ? AND program_id = ?")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = stmt.Exec(hwid, time.Now(), ua.Username, ua.ProgramID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (ua *UserAccount) UpdateLicense(expiry time.Time) {
	stmt, err := DatabaseConnection.Prepare("UPDATE program_users SET expiry = ? WHERE username = ? AND program_id = ?")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = stmt.Exec(expiry, ua.Username, ua.ProgramID)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (ua *UserAccount) GetLastReset() (time.Time, error) {
	stmt, err := DatabaseConnection.Prepare("SELECT last_reset FROM program_users WHERE id = ?")
	if err != nil {
		fmt.Println(err.Error())
		return time.Now(), err
	}

	r := stmt.QueryRow(ua.Id)

	var last time.Time
	err = r.Scan(&last)
	if err == nil {
		return last, nil
	}
	return time.Now(), err
}

func (ua *UserAccount) LinkDiscord(discordid int) error {
	stmt, err := DatabaseConnection.Prepare("INSERT INTO client_discords (user_id, discord_id, program_id) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(ua.Id, discordid, ua.ProgramID)
	return err
}

func (ua *UserAccount) GetDiscord() (string, error) {
	stmt, err := DatabaseConnection.Prepare("SELECT discord_id FROM client_discords WHERE user_id = ?")
	if err != nil {
		log.Fatal(err.Error())
	}

	if r := stmt.QueryRow(ua.Id); r != nil {
		var discId string
		if err = r.Scan(&discId); err == nil {
			return discId, nil
		}
		return "", err
	}
	return "", errors.New("failed to get discord, make sure client is linked")
}

type Token struct {
	Token string
	Level int
	Days int
	Used int
	UsedBy string
	ProgramID int
}

func (t *Token) Use(user *UserAccount) (*time.Time, error) {
	stmt, err := DatabaseConnection.Prepare("UPDATE tokens SET used = 1, used_by = ? WHERE token = ? AND program_id = ?")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(user.Id, t.Token, t.ProgramID)
	if err != nil {
		return nil, err
	}

	_, err = res.RowsAffected()
	if err != nil  {
		return nil, err
	}

	stmt, err = DatabaseConnection.Prepare("UPDATE program_users SET expires = ? WHERE username = ? AND program_id = ?")
	if err != nil {
		return nil, err
	}
	if time.Now().After(user.Expires) {
		user.Expires = time.Now().AddDate(0, 0, t.Days)
	} else {
		user.Expires = user.Expires.AddDate(0, 0, t.Days)
	}

	_, err = stmt.Exec(user.Expires, user.Username, t.ProgramID)
	return &user.Expires, err
}

type Var struct {
	Name string
	Value string
	ProgramID int
}
