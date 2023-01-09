package commands

import (
	"fmt"
	"github.com/xdumpgo/XDG/api/server"
	"github.com/xdumpgo/XDG/apiproto"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/discord"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/utils"
	"github.com/bwmarrin/discordgo"
)

func StatsCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if program, err := server.GetProgram("iMH5RjESHkMJr8EwbcNjYtcsz9fIXEGlp0fXDc5u"); err == nil {
		emb := utils.NewEmbed().SetTitle("Your XDG stats").
			SetThumbnail(discord.BotUser.AvatarURL("250x250")).
			SetColor(utils.RandomColor()).
			SetFooter(utils.FooterTimestamp()).
			SetDescription("Getting stats...")
		msg, err := s.ChannelMessageSendEmbed(m.ChannelID, emb.MessageEmbed)

		user, err := program.GetUserByDiscord(m.Author.ID)
		if err != nil {
			//fmt.Println(err.Error())
			utils.SendError(s, m, "Please make sure you've had your account linked.")
			return
		}

		clientOnline := func() (bool, *server.Client) {
			for _, client := range server.APIServer.GetClients() {
				if client.Name == user.Username {
					return true, client
				}
			}
			return false, nil
		}

		if found, client := clientOnline(); found {
			client.Writer.Write(apiproto.StatsUpdate{})

			status := <- client.Status

			emb.SetDescription(fmt.Sprintf("Good day, %s", user.Username))
			emb.AddField("Current Module", status.CurrentModule)
			emb.AddField("Runtime", status.Runtime.String())
			emb.AddField("Urls", fmt.Sprintf("%d", status.Urls))
			emb.AddField("Injectables", fmt.Sprintf("%d",status.Injectables))
			emb.AddField("Rows", fmt.Sprintf("%d", status.Rows))
			emb.AddField("Threads", fmt.Sprintf("%d", status.Threads))
			emb.AddField("Workers", fmt.Sprintf("%d", status.Workers))
			emb.AddField("Index", fmt.Sprintf("%d / %d", status.Index, status.End))

			emb.InlineAllFields()

			s.ChannelMessageEditEmbed(msg.ChannelID, msg.ID, emb.MessageEmbed)
		} else {
			emb.SetDescription("Looks like your XDG isn't online, you should start it.")
			s.ChannelMessageEditEmbed(msg.ChannelID, msg.ID, emb.MessageEmbed)
			return
		}
	} else {
		fmt.Println(err.Error())
	}
}