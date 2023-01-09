package main

import (
	"database/sql"
	"fmt"
	"github.com/xdumpgo/XDG/commands"
	"github.com/xdumpgo/XDG/config"
	"github.com/xdumpgo/XDG/discord"
	"github.com/xdumpgo/XDG/shoppy"
	"github.com/xdumpgo/XDG/utils"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func init() {

}

func main() {
	utils.DatabaseConnection = CreateCon()
	defer utils.DatabaseConnection.Close()

	go shoppy.WebServerStartup()

	BotObject, err := discordgo.New("Bot " + config.TOKEN)
	if err != nil {
		log.Fatal("Error opening Discord Session")
	}
	defer BotObject.Close()
	discord.BotUser, err = BotObject.User("@me")
	if err != nil {
		log.Fatal("Error getting bot user")
	}

	BotObject.AddHandler(commandHandler)
	BotObject.AddHandler(func(d *discordgo.Session, ready *discordgo.Ready) {
		fmt.Println("Bot Started, running on", len(d.State.Guilds), "servers")
		go func(d *discordgo.Session, ready *discordgo.Ready) {
			for {
				var version string
				if row := utils.DatabaseConnection.QueryRow("SELECT version FROM programs WHERE id = 57"); row != nil {
					row.Scan(&version)
				}

				err = BotObject.UpdateStatus(0, fmt.Sprintf("XDumpGO | v%s ", version))
				time.Sleep(25 * time.Second)
			}
		}(d,ready)
	})

	err = BotObject.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}
	log.Printf("Now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func CreateCon() *sql.DB {
	db, err := sql.Open("mysql", "qauth:fA7FanTBZk^cHLf8@tcp(127.0.0.1:3319)/quartzauth?parseTime=true")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	} else {
		fmt.Println("db is connected")
	}
	//defer db.Close()
	// make sure connection is available
	err = db.Ping()
	fmt.Println(err)
	if err != nil {
		fmt.Println("MySQL db is not connected")
		fmt.Println(err.Error())
	}
	return db
}