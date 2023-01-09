package commands

import (
	"fmt"
	"github.com/xdumpgo/XDG/api/server"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/discord"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/utils"
	"github.com/bwmarrin/discordgo"
	"time"
)

func DownloadCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if program, err := server.GetProgram("iMH5RjESHkMJr8EwbcNjYtcsz9fIXEGlp0fXDc5u"); err == nil {
		emb := utils.NewEmbed().SetTitle("Get XDumpGO Files").
			SetThumbnail(discord.BotUser.AvatarURL("250x250")).
			SetColor(0x00ff00).
			SetFooter(utils.FooterTimestamp())

		//		discId := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(args[1], "<@", ""), ">", ""), "!", "")

		if user, err := program.GetUserByDiscord(m.Author.ID); err == nil {
			if user.Expires.After(time.Now()) {
				if dm, err := s.UserChannelCreate(m.Author.ID); err == nil {
					emb.SetDescription(fmt.Sprintf("%s, Sent you the download link in your DM's!", m.Author.Mention()))
					s.ChannelMessageSend(dm.ID, fmt.Sprintf("Here's the download for XDumpGO v%s - [%s]\nYour username is `%s`", program.Version, program.Url.String, user.Username))
				} else {
					emb.SetColor(0xff0000)
					emb.SetDescription(fmt.Sprintf("%s, Please open your DM's to receive the download link!", m.Author.Mention()))
				}
			} else {
				emb.SetColor(0xff0000)
				emb.SetDescription(fmt.Sprintf("%s, You do not have an active license.  Please buy a new license @ https://quartzinc.dev/xdumpgo", m.Author.Mention()))
			}
		} else {
			//utils.SendError(s,m, "Please have your account linked by an administrator")
			emb.SetColor(0xff0000)
			emb.SetDescription("Please have your account linked by an administrator")
		}
		s.ChannelMessageSendEmbed(m.ChannelID, emb.MessageEmbed)
	} else {
		fmt.Println(err.Error())
	}
}