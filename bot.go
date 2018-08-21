package main

import (
	"github.com/IoTeXGroupIT/IoTeXGroupITBot/botApi"
	"github.com/IoTeXGroupIT/IoTeXGroupITBot/configuration"
	"github.com/IoTeXGroupIT/IoTeXGroupITBot/reminder"
	"github.com/IoTeXGroupIT/IoTeXGroupITBot/spamFilter"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func main() {
	log.Printf("[START] Starting the bot...")

	API := configuration.GetApiKey()
	log.Printf("Using Telegram API Key %s", API)

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
			case "events":
				reminder.List(bot, chat)
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

				// Check spam messages into new user first/last name and kick it
				if !spamFilter.FilterNewUserJoined(bot, message) {
					continue
				}

				// If everything is good send welcome message
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
