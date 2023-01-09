package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/xdumpgo/XDG/api/server"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/commands"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/config"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/discord"
	"github.com/xdumpgo/XDG/cmd/apiserver/discord/utils"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"strings"
	"time"
)

func main() {
	address := flag.String("address", "0.0.0.0:7005", "Set the listening address.")

	flag.Parse()

	server.DatabaseConnection = CreateCon()

	StartDiscordBot()

	server.APIServer = server.NewServer()
	err := server.APIServer.Listen(*address)
	if err != nil {
		log.Fatal(err.Error())
	}

	// start the server
	server.APIServer.Start()
}

func StartDiscordBot() {
	commands.Commands = new(commands.CommandList)
	fmt.Println("Adding discord bot commands")
	commands.Commands.NewCommand("help", []string{"help", "?"}, "Shows this help message", "help [cmd]", 0, discordgo.PermissionSendMessages, []string{"cmd"}, commands.HelpCommand)
	commands.Commands.NewCommand("usage", []string{"usage", "u"}, "Shows command usage", "usage [cmd]", 0, discordgo.PermissionSendMessages, []string{"cmd"}, commands.UsageCommand)
	//commands.Commands.NewCommand("ping", []string{"ping", "pi"}, "Pings the bot", "ping", 0, discordgo.PermissionSendMessages, []string{}, PingCommand)
	commands.Commands.NewCommand("plan", []string{"plan", "pl"}, "Displays your plan", "plan <@user>", 0, discordgo.PermissionSendMessages, []string{"mention"}, commands.PlanCommand)
	commands.Commands.NewCommand("link", []string{"link", "l"}, "Links a XDG user to a discord account", "link [username] [@user]", 2, discordgo.PermissionBanMembers, []string{"user", "mention"}, commands.LinkCommand)
	commands.Commands.NewCommand("hwid", []string{"hwid", "hw"}, "Mount your HWID for XDumpGO", "hwid <user/@user>", 0, discordgo.PermissionSendMessages, []string{}, commands.HWIDCommand)
	commands.Commands.NewCommand("redeem", []string{"redeem", "r"}, "Redeem a token to extend your license", "redeem [token]", 1, discordgo.PermissionSendMessages, []string{"token"}, commands.RedeemCommand)
	commands.Commands.NewCommand("download", []string{"download", "d"}, "Get a link to latest version of XDumpGO", "download", 0, discordgo.PermissionSendMessages, []string{}, commands.DownloadCommand)
	commands.Commands.NewCommand("stats", []string{"stats", "s"}, "Display stats of your XDG instance", "stats", 0, discordgo.PermissionSendMessages, []string{}, commands.StatsCommand)
	commands.Commands.NewCommand("users", []string{"users", "u"}, "Shows active users", "users", 0, discordgo.PermissionSendMessages, []string{}, commands.UserListCommand)
	commands.Commands.NewCommand("term", []string{"term", "t"}, "Terminate active user session", "term [username] [reason]", 1, discordgo.PermissionBanMembers, []string{"username", "reason"}, commands.TermCommand)
	commands.Commands.NewCommand("melt", []string{"melt", "del"}, "Melt's a users files.", "melt [username]", 1, discordgo.PermissionBanMembers, []string{"username"}, commands.MeltCommand)

	var err error
	discord.BotObject, err = discordgo.New("Bot " + config.TOKEN)
	if err != nil {
		log.Fatal("Error opening Discord Session")
	}
	defer discord.BotObject.Close()
	discord.BotUser, err = discord.BotObject.User("@me")
	if err != nil {
		log.Fatal("Error getting bot user")
	}

	discord.BotObject.AddHandler(commandHandler)
	discord.BotObject.AddHandler(func(d *discordgo.Session, ready *discordgo.Ready) {
		fmt.Println("Bot Started, running user", discord.BotUser.Username, "on", len(d.State.Guilds), "servers")
		go func(d *discordgo.Session, ready *discordgo.Ready) {
			for {
				var version string
				if row := server.DatabaseConnection.QueryRow("SELECT version FROM programs WHERE id = 57"); row != nil {
					row.Scan(&version)
				}

				err = discord.BotObject.UpdateStatus(0, fmt.Sprintf("XDumpGO | v%s ", version))
				time.Sleep(25 * time.Second)
			}
		}(d,ready)
		go func(s *discordgo.Session) {
			for {
				select {
				case <- time.After(1 * time.Hour):
					if rows, err := server.DatabaseConnection.Query("SELECT id FROM  program_users WHERE program_id=57 AND expiry < NOW()"); err == nil {
						var userIds []int
						var id int
						for rows.Next() {
							if rows.Scan(&id) == nil {
								userIds = append(userIds, id)
							}
						}

						var did int
						for _, uid := range userIds {
							row := server.DatabaseConnection.QueryRow("SELECT discord_id FROM client_discords WHERE user_id = ?", uid)
							if row.Scan(&did) == nil {
								s.GuildMemberRoleRemove("748271210676355153", strconv.Itoa(did), "748272709548376214")
							}
						}
					}
				}
			}
		}(d)
	})

	err = discord.BotObject.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}

	fmt.Println("Discord bot started")
}

func CreateCon() *sql.DB {
	db, err := sql.Open("mysql", "qauth:fA7FanTBZk^cHLf8@tcp(localhost:3306)/quartzauth?parseTime=true")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	} else {
		fmt.Println("db is connected")
	}
	//defer db.Close()
	// make sure connection is available
	if err = db.Ping(); err != nil {
		fmt.Println("MySQL db is not connected")
		fmt.Println(err.Error())
	}
	return db
}

func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	cmd := m.Content

	// Check for command prefix


	if strings.HasPrefix(cmd, config.CMD_PREFIX) {
	//	if c,err := s.Channel(m.ChannelID); err != nil && c != nil {
			//if strings.Contains(c.Name, "bot-spam") || strings.Contains(c.Name, "ticket-") {
				cmd = strings.TrimPrefix(cmd, config.CMD_PREFIX)
				parts := strings.Split(cmd, " ")
				cmd = parts[0]
				args := parts[1:]
				if command, ok := commands.Commands.Contains(cmd); ok {
					if len(args) < command.MinArgs {
						go utils.SelfDestructingMessage(s,m, "Error", "Not enough arguments, see help", 5)
					} else {
						command.Execute(s, m, args)
					}
				}
			//} else {
			//	go utils.SelfDestructingMessage(s,m, "Error", "Please use commands in bot-spam or your ticket.", 5)
			//}
		//} else {
	//		go utils.SelfDestructingMessage(s,m, "Error", "Please use commands in bot-spam or your ticket.", 5)
	//	}
	}
}