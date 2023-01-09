package shoppy

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

const secret = ""

func WebServerStartup() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)

	log.Fatal(http.ListenAndServe(":6797", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	raw, _ := ioutil.ReadAll(r.Body)
	ioutil.WriteFile("debug.txt", raw, 0644)
	/*raw, _ := ioutil.ReadAll(r.Body)
	var aa ShoppyOrder
	err := json.Unmarshal(raw, &aa)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(aa)
	a := aa.Data.Order.CustomFields[0].Value

	if !strings.Contains(a, "#") {
		emb := utils.NewEmbed().SetTitle("Automatic Support Request").SetDescription("User didn't provide discriminator, manual rank upgrade required. \nUsername: " + a).SetColor(0x33a5ff )
		discord.BotObject.ChannelMessageSendEmbed("685707946876731561", emb.MessageEmbed)
	} else {
		if verifyMessage(raw, r.Header.Get("X-Shoppy-Signature")) {
			for _, guild := range discord.BotObject.State.Guilds {
				roles, _ := discord.BotObject.GuildRoles(guild.ID)
				for _, member := range guild.Members {
					z := strings.Split(a, "#")
					if member.User.Username == z[0] && member.User.Discriminator == z[1] {
						for _, n := range roles {
							if n.Name == "VIP" {
								err = discord.BotObject.GuildMemberRoleAdd(member.GuildID, member.User.ID, n.ID)
								fmt.Println("Added vip to member " + member.User.Username + " email: " + aa.Data.Order.Email)
								emb := utils.NewEmbed().SetTitle("Automatic Rank Upgrade").
									SetColor(0x3cff33).
									SetDescription(fmt.Sprintf("User: %s\nThank you for your support!  Your rank has been upgraded!", member.Mention()))
								discord.BotObject.ChannelMessageSendEmbed("685707946876731561", emb.MessageEmbed)
								if err != nil {
									fmt.Println("poop " + err.Error())
								}
							}
						}
					}
				}
			}
		} else {
			fmt.Println("error verifying")
		}
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(""))
	if err != nil {
		fmt.Println(err.Error())
	}*/
}

func verifyMessage(message []byte, signature string) bool {
	/*tmp := ComputeHmac256(message, []byte(secret))
	//fmt.Println(tmp)
	if tmp == signature {
		return true
	}*/
	return false
}