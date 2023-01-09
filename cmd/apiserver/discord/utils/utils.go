package utils

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/AvraamMavridis/randomcolor"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mitchellh/go-ps"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)


func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, 14)
	if err != nil {
		log.Println(err)
	}    // GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func FmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

type Field struct {
	name string
	value string
}

func SendError(s *discordgo.Session, m* discordgo.MessageCreate, message string) {
	emb := NewEmbed().
		SetTitle("ERROR!").
		SetDescription(message).
		SetColor(0xff0000)
	s.ChannelMessageSendEmbed(m.ChannelID, emb.MessageEmbed)
}

func ArrayContains(arr []string, str string) bool {
	for _, k := range arr {
		if k == str {
			return true
		}
	}
	return false
}

func ArrayReverse(a []string) []string {
	b := a
	for i := len(b)/2-1; i >= 0; i-- {
		opp := len(b)-1-i
		b[i], b[opp] = b[opp], b[i]
	}
	return b
}

func UserHasRole(s *discordgo.Session, m *discordgo.MessageCreate, roleIDInput string) bool {
	member, err := s.State.Member(m.GuildID, m.Author.ID)
	if err != nil {
		if member, err = s.GuildMember(m.GuildID, m.Author.ID); err != nil {
			SendError(s,m, "errr.")
			return false
		}
	}
	for _, roleID := range member.Roles {
		if roleID == roleIDInput {
			return true
		}
	}
	return false
}

func ComputeHmac256(message []byte, secret []byte) string {
	hash := hmac.New(sha512.New, secret)
	hash.Write(message)

	// to lowercase hexits
	return hex.EncodeToString(hash.Sum(nil))

	//return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func SelfDestructingMessage(s *discordgo.Session, m *discordgo.MessageCreate, title string, message string, secondsToDeath int) {
	emb := NewEmbed().SetTitle(title).SetDescription(message).SetColor(0xff9333)
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		fmt.Println(err.Error())
	}
	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, emb.MessageEmbed)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for i:=secondsToDeath; i>0; i-- {
		emb.SetFooter(fmt.Sprintf("Seconds to self destruct: %d", i))
		msg, _ = s.ChannelMessageEditEmbed(msg.ChannelID, msg.ID, emb.MessageEmbed)
		delaySecond(1)
	}
	s.ChannelMessageDelete(msg.ChannelID, msg.ID)
}

func delaySecond(n time.Duration) {
	time.Sleep(n * time.Second)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrc(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func FooterTimestamp() string {
	return fmt.Sprintf("%s", time.Now().Format("Mon Jan 2 15:04:05"))
}

func RandomColor() int {
	color, err := strconv.ParseInt(strings.ReplaceAll(randomcolor.GetRandomColorInHex(), "#", ""), 16, 64)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	return int(color)
}

func ComesFromDM(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		if channel, err = s.Channel(m.ChannelID); err != nil {
			return false, err
		}
	}

	return channel.Type == discordgo.ChannelTypeDM, nil
}

func CreationTime(ID string) (t time.Time, err error) {
	i, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return
	}
	timestamp := (i >> 22) + 1420070400000
	t = time.Unix(timestamp/1000, 0)
	return
}

func FindProcess(name string) bool {
	procs, err := ps.Processes()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	for _,k := range procs {
		if k.Executable() == name {
			return true
		}
	}
	return false
}