package commands

import (
	"fmt"
	"github.com/xdumpgo/XDG/api/server"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/discord"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/utils"
	"github.com/bwmarrin/discordgo"
)

func RedeemCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if program, err := server.GetProgram("iMH5RjESHkMJr8EwbcNjYtcsz9fIXEGlp0fXDc5u"); err == nil {
		emb := utils.NewEmbed().SetTitle("Redeem a token").
			SetThumbnail(discord.BotUser.AvatarURL("250x250")).
			SetColor(utils.RandomColor()).
			SetFooter(utils.FooterTimestamp())

		if user, err := program.GetUserByDiscord(m.Author.ID); err == nil {
			if token, err := program.GetToken(args[0]); err == nil {
				 if expiry, err := token.Use(user); err == nil {
				 	emb.SetDescription(fmt.Sprintf("Used token, your license will now expire on `%s`", expiry.Format("Mon Jan _2 15:04:05 2006")))
				 } else {
				 	utils.SendError(s,m, "Failed to redeem token, contact an administrator")
				 	return
				 }
			} else {
				fmt.Println(err.Error())
				emb.SetDescription("Invalid token.")
			}
			s.ChannelMessageSendEmbed(m.ChannelID, emb.MessageEmbed)
		} else {
			utils.SendError(s,m, "Please have your account linked by an administrator")
		}

	} else {
		fmt.Println(err.Error())
	}
}