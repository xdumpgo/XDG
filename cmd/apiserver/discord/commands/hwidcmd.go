package commands

import (
	"fmt"
	"github.com/xdumpgo/XDG/api/server"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/discord"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/utils"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

func HWIDCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if program, err := server.GetProgram("iMH5RjESHkMJr8EwbcNjYtcsz9fIXEGlp0fXDc5u"); err == nil {
		emb := utils.NewEmbed().SetTitle("HWID Mount").
			SetThumbnail(discord.BotUser.AvatarURL("250x250")).
			SetColor(0x00ff00).
			SetFooter(utils.FooterTimestamp())

		if len(args) > 0 && UserHasPermission(s,m, "iMH5RjESHkMJr8EwbcNjYtcsz9fIXEGlp0fXDc5u") {
			if user, err := program.GetUserByDiscord(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(args[0], "<@", ""), ">", ""), "!", "")); err == nil {
				user.SetHWID("")
				emb.SetDescription(fmt.Sprintf("Successfully reset HWID for user %s.", user.Username))
			} else {
				emb.SetDescription("Unknown user").SetColor(0xff0000)
			}
		} else {
			if user, err := program.GetUserByDiscord(m.Author.ID); err == nil {
				if last, err := user.GetLastReset(); err == nil {
					if last.AddDate(0,0,7).Before(time.Now()) {
						user.SetHWID("")
						emb.SetDescription("Successfully reset HWID.")
					} else {
						emb.SetColor(0xff0000)
						emb.SetDescription(fmt.Sprintf("Failed to reset HWID, your last reset was on `%s`, you must wait until `%s` to reset again.", user.LastReset.Format("Mon Jan _2 15:04:05 2006"), user.LastReset.AddDate(0, 0, 7).Format("Mon Jan _2 15:04:05 2006")))
					}
				}
			} else {
				emb.SetDescription("Please have your account linked by an administrator")
				emb.SetColor(0xff0000)
			}
		}
		s.ChannelMessageSendEmbed(m.ChannelID, emb.MessageEmbed)
	} else {
		fmt.Println(err.Error())
	}
}