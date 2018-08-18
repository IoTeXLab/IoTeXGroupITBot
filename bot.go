package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"github.com/simonerom/IoTeXGroupITBot/configuration"
	"github.com/simonerom/IoTeXGroupITBot/botApi"
)


func main() {
    log.Printf("[START] Starting the bot...")

	API := configuration.GetApiKey();

	bot, err := tgbotapi.NewBotAPI(API)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {

		// Ignore null messages
		if update.Message == nil {
			continue
		}

		message := update.Message
		chat := message.Chat

		// Process commands
		if message.IsCommand() {
			switch message.Command() {

			case "help": // Shows available commands
				botApi.PostHelpMessage(bot, chat, message)
			case "roadmap": // Post the roadmap image
				botApi.PostRoadmapImage(bot, chat, message)
			default:
				botApi.PostWrongCommandMessage(bot, chat, message)
			}
			continue
		}

		// Manage new users joined messages
		if message.NewChatMembers != nil {

			for _, user := range *message.NewChatMembers {
				// We log any new join
				LogNewUserJoined(message.Chat, user)

				// Send welcome message
				if configuration.Cfg.PostWelcomeMessage {
					botApi.PostWelcomeMessage(bot, chat, message)
				}
			}
		}

	}
}



// Log the new user join event
func LogNewUserJoined(chat *tgbotapi.Chat, user tgbotapi.User) {
	log.Printf("_________________________")
	log.Printf("New user joined the group %s", chat.UserName)
	log.Printf("Username: %s", user.UserName)
	log.Printf("First Name: %s", user.FirstName)
	log.Printf("Last Name: %s", user.LastName)
}

// Obtain the Telegram Bot Api key either by environment variable or configuration.json

