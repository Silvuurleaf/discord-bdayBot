package bot

import (
	"discordBot/config"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var BotID string
var goBot *discordgo.Session

const prefix string = "!bdaybot"

type Answers struct {
	OriginChannelId string
	Birthday        string
}

func (a *Answers) ToMessageEmbed() discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "Birthday",
			Value: a.Birthday,
		},
	}

	return discordgo.MessageEmbed{
		Title:  "New Responses!",
		Fields: fields,
	}
}

var responses map[string]Answers = map[string]Answers{}

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	goBot.AddHandler(messageHandler)

	goBot.AddHandler(userPromptHandler)

	goBot.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer goBot.Close()

	fmt.Println("Bot is running!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		return
	}

	//DM lgoic

	if m.GuildID == "" {
		answers, ok := responses[m.ChannelID]

		if !ok {
			return
		}
		if answers.Birthday == "" {
			answers.Birthday = m.Content

			s.ChannelMessageSend(m.ChannelID, "Great! Thanks")
			responses[m.ChannelID] = answers

			return
		} else {
			embed := answers.ToMessageEmbed()
			s.ChannelMessageSendEmbed(answers.OriginChannelId, &embed)
			delete(responses, m.ChannelID)
		}
	}

	args := strings.Split(m.Content, " ")

	if args[0] != prefix {
		return
	}

	if args[1] == "prompt" {
		userPromptHandler(s, m)
	}

	// If the message is "Hi" reply with "Hi Back!!"
	if args[1] == "Hi" {

		embed := discordgo.MessageEmbed{
			Title: "Hi Back",
			URL:   "https://www.youtube.com/",
		}

		s.ChannelMessageSendEmbed(m.ChannelID, &embed)

		_, _ = s.ChannelMessageSend(m.ChannelID, "Hi Back")
	}

}

func userPromptHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	//if user answered questions ignore otherwise ask

	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		log.Panic(err)
	}

	if _, ok := responses[channel.ID]; !ok {
		responses[channel.ID] = Answers{
			OriginChannelId: m.ChannelID,
			Birthday:        "",
		}

		s.ChannelMessageSend(channel.ID, "What's your Birthday?")
	} else {
		s.ChannelMessageSend(channel.ID, "Please respond to me );=")
	}

	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		return
	}

	args := strings.Split(m.Content, " ")

	if args[0] != prefix {
		return
	}

}
