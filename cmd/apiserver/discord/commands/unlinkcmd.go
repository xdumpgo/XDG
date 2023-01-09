package commands

import (
	"fmt"
	"github.com/xdumpgo/XDG/api/server"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/discord"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/utils"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

func UnlinkCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !UserHasPermission(s,m, "iMH5RjESHkMJr8EwbcNjYtcsz9fIXEGlp0fXDc5u") {
		utils.SendError(s,m, "Sorry, this is only available for Staff members.")
		return
	}

	if program, err := server.GetProgram("iMH5RjESHkMJr8EwbcNjYtcsz9fIXEGlp0fXDc5u"); err == nil {
		emb := utils.NewEmbed().SetTitle(fmt.Sprintf("Unlinking user account for %s", program.Name)).
			SetThumbnail(discord.BotUser.AvatarURL("250x250")).
			SetColor(utils.RandomColor()).
			SetFooter(utils.FooterTimestamp())
		if len(args) == 2 {
			if user, err := program.GetUserByName(args[0]); err == nil {
				a := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(args[1], "<@", ""), ">", ""), "!", "")
				did, _ := strconv.Atoi(a)
				if err := user.LinkDiscord(did); err != nil {
					fmt.Println(err.Error())
					utils.SendError(s,m, "Failed to link account")
					return
				}
				emb.SetDescription(fmt.Sprintf("Successfully linked <@%d> to XDG user account %s", did, user.Username))
				s.GuildMemberRoleAdd(m.GuildID, a, "710735787322507344")
			} else {
				fmt.Println(err.Error())
				utils.SendError(s, m, "Invalid XDG user account.")
				return
			}
			s.ChannelMessageSendEmbed(m.ChannelID, emb.MessageEmbed)
		} else {
			utils.SendError(s,m, "Please see command usage")
		}
	} else {
		fmt.Println(err.Error())
	}
}