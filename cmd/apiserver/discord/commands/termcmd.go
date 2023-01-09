package commands

import (
	"fmt"
	"github.com/xdumpgo/XDG/api/server"
	protocol "github.com/xdumpgo/XDG/apiproto"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/discord"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/utils"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func TermCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !UserHasPermission(s,m, "iMH5RjESHkMJr8EwbcNjYtcsz9fIXEGlp0fXDc5u") {
		utils.SendError(s,m, "Sorry, this is only available for Staff members.")
		return
	}

	if _, err := server.GetProgram("iMH5RjESHkMJr8EwbcNjYtcsz9fIXEGlp0fXDc5u"); err == nil {
		emb := utils.NewEmbed().SetTitle(fmt.Sprintf("Attempting to terminate user %s", args[0])).
			SetThumbnail(discord.BotUser.AvatarURL("250x250")).
			SetColor(utils.RandomColor()).
			SetFooter(utils.FooterTimestamp())
		emb.SetDescription("Failed to find user")

		for _, client := range server.APIServer.GetClients() {
			if client.Name == args[0] {
				emb.SetDescription(fmt.Sprintf("Terminated user session for %s", args[0]))
				client.Writer.Write(protocol.Terminate{Reason: strings.Join(args[1:], " ")})
			}
		}

		s.ChannelMessageSendEmbed(m.ChannelID, emb.MessageEmbed)
	} else {
		fmt.Println(err.Error())
	}
}