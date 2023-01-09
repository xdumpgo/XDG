package commands

import (
	"github.com/xdumpgo/XDG/api/server"
	"github.com/bwmarrin/discordgo"
)

var SellixAuth = "iuHCiRBTxBIiUSewhAl6AdPX0Bts4JIkklUsLiDYwU9qVFNuAo4gCOWfDmD60nn1"
var SellixEndpoint = "https://dev.sellix.io/v1"

type Command struct {
	Ident string
	Commands []string
	MinArgs int
	Args []string
	Help string
	Usage string
	Permission int
	Execute func(*discordgo.Session, *discordgo.MessageCreate, []string)
}

type CommandList struct {
	commands []*Command
}

var Commands *CommandList

func (list* CommandList) NewCommand(ident string, cmd []string, help string, usage string, minArgs int, permission int, args []string, execute func(*discordgo.Session, *discordgo.MessageCreate, []string)) {
	tmp := &Command{Ident: ident, Commands:cmd, Help:help, Usage: usage, MinArgs:minArgs, Permission:permission, Args:args, Execute: execute}
	list.commands = append(list.commands, tmp)
}

func (list* CommandList) Contains(cmd string) (*Command, bool) {
	for _, k := range list.commands {
		if k.Contains(cmd) {
			return k, true
		}
	}
	return nil, false
}

func (c *Command) Contains(a string) bool {
	for _, k := range c.Commands {
		if k == a {
			return true
		}
	}
	return false
}

func UserHasPermission (s *discordgo.Session, m *discordgo.MessageCreate, key string) bool {
	for _, program := range server.GetProgramsByUser(server.GetFrontendLinkUser(m.Author.ID)) {
		if program.Key == key {
			return true
		}
	}
	return false
}

func MemberHasPermission(s *discordgo.Session, guildID string, userID string, permission int) (bool, error) {
	member, err := s.State.Member(guildID, userID)
	if err != nil {
		if member, err = s.GuildMember(guildID, userID); err != nil {
			return false, err
		}
	}

	// Iterate through the role IDs stored in member.Roles
	// to check permissions
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			return false, err
		}
		if role.Permissions&permission != 0 {
			return true, nil
		}
	}

	return false, nil
}