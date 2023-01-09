package commands

import (
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/utils"
	"github.com/bwmarrin/discordgo"
)

func ValidateOrderCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !UserHasPermission(s,m, "iMH5RjESHkMJr8EwbcNjYtcsz9fIXEGlp0fXDc5u") {
		utils.SendError(s,m, "Sorry, this is only available for Staff members.")
		return
	}

	if len(args) != 1 {
		utils.SendError(s,m, "Invalid arguments, see help.")
		return
	}
/*
	emb := utils.NewEmbed().SetTitle(fmt.Sprintf("List of active users for %s", program.Name)).
		SetThumbnail(discord.BotUser.AvatarURL("250x250")).
		SetColor(utils.RandomColor()).
		SetFooter(utils.FooterTimestamp())
*/

}