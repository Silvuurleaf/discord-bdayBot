package bot

import "discordBot/config"

var BotID string
var goBot *discordgo.Session

func Start() {
	discordgo.New("Bot" + config.Token)
}
