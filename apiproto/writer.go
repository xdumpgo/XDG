package apiproto

import (
	"encoding/json"
	"github.com/xdumpgo/XDG/utils"
	"github.com/oreans/virtualizersdk"
	"io"
)

type CommandWriter struct {
	writer io.Writer
}

func NewCommandWriter(writer io.Writer) *CommandWriter {
	return &CommandWriter{
		writer: writer,
	}
}

func (w *CommandWriter) writeString(command byte, msg string) error {
	virtualizersdk.Macro(virtualizersdk.SHARK_WHITE_START)
	payload := []byte{command}
	payload = append(payload, []byte(utils.Encrypt(msg, "zLy&CyLg#tvUp4aH6tDH7F%fx=w6xCa%r_A^3hkq468nuuc5=7xba^un&bYLeQ=C-qL_hp#rnNU5a!5neb%_&aygXDL8Jg7?Y6cHpLAtEXM&kbfGTzyQm3Lv"))...)
	payload = append(payload, '\n')
	//fmt.Printf("sendv [%#x] %s l:%d\n", command, msg, len(payload))

	_, err := w.writer.Write(payload)

	virtualizersdk.Macro(virtualizersdk.SHARK_WHITE_END)

	return err
}

func (w *CommandWriter) Write(command interface{}) error {
	// naive implementation ...
	var err error

	raw, err := json.Marshal(command)
	if err != nil {
		return err
	}

	switch command.(type) {
	case Heartbeat:
		err = w.writeString(HEARTBEAT, string(raw))
	case SendCommand:
   		err = w.writeString(SEND, string(raw))
	case MessageCommand:
		err = w.writeString(MESSAGE, string(raw))
	case UserList:
		err = w.writeString(USERLIST, string(raw))
	case LoginCommand:
		err = w.writeString(LOGIN, string(raw))
	case RegisterCommand:
		err = w.writeString(REGISTER, string(raw))
	case RedeemCommand:
		err = w.writeString(REDEEM, string(raw))
	case VarCommand:
		err = w.writeString(VAR, string(raw))
	case StatsUpdate:
		err = w.writeString(STATUS, string(raw))
	case AuthResponse:
		err = w.writeString(RESPONSE, string(raw))
	case Terminate:
		err = w.writeString(TERMINATE, string(raw))
	case Melt:
		err = w.writeString(MELT, string(raw))
	case Proxies:
		err = w.writeString(PROXYLIST, string(raw))
	case NewsCommand:
		err = w.writeString(NEWS, string(raw))
	case NameCommand:
		err = w.writeString(NAME, string(raw))
	default:
		err = UnknownCommand
	}

	return err
}