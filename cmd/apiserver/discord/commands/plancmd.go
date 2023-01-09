package commands

import (
	"fmt"
	"github.com/xdumpgo/XDG/api/server"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/discord"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/utils"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func PlanCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if program, err := server.GetProgram("iMH5RjESHkMJr8EwbcNjYtcsz9fIXEGlp0fXDc5u"); err == nil {
		emb := utils.NewEmbed().SetTitle(fmt.Sprintf("Your plan for %s", program.Name)).
			SetThumbnail(discord.BotUser.AvatarURL("250x250")).
			SetColor(utils.RandomColor()).
			SetFooter(utils.FooterTimestamp())
		if len(args) > 0 && server.GetFrontendLinkUser(m.Author.ID) != 0 {
			if user, err := program.GetUserByDiscord(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(args[0], "<@", ""), ">", ""), "!", "")); err == nil {
				emb.SetDescription(fmt.Sprintf("Good day, %s", user.Username))
				emb.AddField("License Expires On", user.Expires.Format("Mon Jan _2 15:04:05 2006"))
			} else {
				utils.SendError(s, m, "Discord not linked to XDG user account.")
			}
		} else {
			user, err := program.GetUserByDiscord(m.Author.ID)
			if err != nil {
				fmt.Println(err.Error())
				utils.SendError(s, m, "Please make sure you've had your account linked.")
				return
			}
			emb.SetDescription(fmt.Sprintf("Good day, %s", user.Username))
			emb.AddField("License Expires On", user.Expires.Format("Mon Jan _2 15:04:05 2006"))
		}

		s.ChannelMessageSendEmbed(m.ChannelID, emb.MessageEmbed)
	} else {
		fmt.Println(err.Error())
	}
}