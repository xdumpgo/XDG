package apiproto

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/xdumpgo/XDG/utils"
	"github.com/oreans/virtualizersdk"
	"io"
	"log"
)

type CommandReader struct {
	reader *bufio.Reader
}

func NewCommandReader(reader io.Reader) *CommandReader {
	return &CommandReader{
		reader: bufio.NewReader(reader),
	}
}

func (r *CommandReader) Read() (interface{}, error) {
	// Read the first part
	virtualizersdk.Macro(virtualizersdk.SHARK_WHITE_START)
	commandName, err := r.reader.ReadByte()

	if err != nil {
		fmt.Print(err.Error() + "\n")
		return nil, err
	}

	payload, err := r.reader.ReadBytes(byte('\n'))
	if err != nil {
		fmt.Print(err.Error() + "\n")
		return nil, err
	}

	b := utils.Decrypt(string(payload), "zLy&CyLg#tvUp4aH6tDH7F%fx=w6xCa%r_A^3hkq468nuuc5=7xba^un&bYLeQ=C-qL_hp#rnNU5a!5neb%_&aygXDL8Jg7?Y6cHpLAtEXM&kbfGTzyQm3Lv")

	switch commandName {
	case HEARTBEAT:
		var msg Heartbeat
		if err = json.Unmarshal([]byte(b), &msg); err == nil {
			return msg, nil
		}
	case MESSAGE:
		var msg MessageCommand
		if err = json.Unmarshal([]byte(b), &msg); err == nil {
			return msg, nil
		}
	case SEND:
		var msg SendCommand
		if err = json.Unmarshal([]byte(b), &msg); err == nil {
			return msg, nil
		}
	case USERLIST:
		var msg UserList
		if err = json.Unmarshal([]byte(b), &msg); err == nil {
			return msg, nil
		}
	case LOGIN:
		var msg LoginCommand
		if err = json.Unmarshal([]byte(b), &msg); err == nil {
			return msg, nil
		}
	case REGISTER:
		var msg RegisterCommand
		if err = json.Unmarshal([]byte(b), &msg); err == nil {
			return msg, nil
		}
	case REDEEM:
		var msg RedeemCommand
		if err = json.Unmarshal([]byte(b), &msg); err == nil {
			return msg, nil
		}
	case VAR:
		var msg VarCommand
		if err = json.Unmarshal([]byte(b), &msg); err == nil {
			return msg, nil
		}
	case RESPONSE:
		var msg AuthResponse
		if err = json.Unmarshal([]byte(b), &msg); err == nil {
			return msg, nil
		}
	case STATUS:
		var msg StatsUpdate
		if err = json.Unmarshal([]byte(b), &msg); err == nil {
			return msg, nil
		}
	case NEWS:
		var msg NewsCommand
		if err = json.Unmarshal([]byte(b), &msg); err == nil {
			return msg, nil
		}
	case TERMINATE:
		var msg Terminate
		if err = json.Unmarshal([]byte(b), &msg); err == nil {
			return msg, nil
		}
	case MELT:
		var msg Melt
		if err = json.Unmarshal([]byte(b), &msg); err == nil {
			return msg, nil
		}
	case PROXYLIST:
		var msg Proxies
		if err = json.Unmarshal([]byte(b), &msg); err == nil {
			return msg, nil
		}
	case NAME:
		var msg NameCommand
		if err = json.Unmarshal([]byte(b), &msg); err == nil {
			return msg, nil
		}
	default:
		log.Printf("Unknown command: %v", commandName)
	}

	virtualizersdk.Macro(virtualizersdk.SHARK_WHITE_END)
	return nil, UnknownCommand
}

func (r *CommandReader) ReadAll() ([]interface{}, error) {
	commands := []interface{}{}

	for {
		command, err := r.Read()

		if command != nil {
			commands = append(commands, command)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return commands, err
		}
	}

	return commands, nil
}