package commands

import (
	"fmt"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/config"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/discord"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/utils"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func UsageCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) > 0 {
		if cmd, ok := Commands.Contains(args[0]); ok {
			printUsageForCommand(s,m,cmd)
		} else {
			utils.SendError(s,m,"Invalid command")
		}
	} else {
		emb := utils.NewEmbed().
			SetTitle("Command Usage").
			SetDescription("A list of commands this bot can handle").
			SetThumbnail(discord.BotUser.AvatarURL("250x250")).
			SetColor(utils.RandomColor()).
			SetFooter(utils.FooterTimestamp())
		for _, k := range Commands.commands {
			emb.AddField(strings.Join(k.Commands, " | "), config.CMD_PREFIX + k.Usage)
		}

		_, err := s.ChannelMessageSendEmbed(m.ChannelID, emb.MessageEmbed)
		if err != nil {
			fmt.Println("ERR!\n" + err.Error())
		}
	}
}

func printUsageForCommand(s *discordgo.Session, m *discordgo.MessageCreate, cmd *Command) {
	emb := utils.NewEmbed().
		SetTitle(cmd.Ident).
		SetDescription(cmd.Usage).
		SetThumbnail(discord.BotUser.AvatarURL("250x250")).
		SetColor(utils.RandomColor()).
		SetFooter(utils.FooterTimestamp())
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, emb.MessageEmbed)
	if err != nil {
		fmt.Println("ERR!\n" + err.Error())
	}
}