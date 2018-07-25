package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	ApiKey string
	KickOnFirstNameLength bool
	FirstNameMaxLength int
	KickOnFullNameLength bool
	FullNameMaxLength int
	CommandsWhitelist []string
}

var cfg Configuration

var bot *tgbotapi.BotAPI

func main() {

	err := gonfig.GetConf("configuration.json", &cfg)

	if err != nil {
		log.Fatal(err)
	}

	API := GetApiKey();

	bot, err = tgbotapi.NewBotAPI(API)
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
				PostHelp(chat, message)
			case "roadmap": // Post the roadmap image
				PostRoadmap(chat, message)
			}
			continue
		}

		// Kicking new users when they include long text in First/Last name fields
		if message.NewChatMembers != nil {

			for _, user := range *message.NewChatMembers {

				LogNewUserJoined(message.Chat, user)

				firstNameLength := len(user.FirstName)
				fullNameLength := firstNameLength + len(user.LastName)
				// TODO: if we want sto stay even safer, whe could kick only chinese text?

				if cfg.KickOnFirstNameLength && firstNameLength > cfg.FirstNameMaxLength {
					// delete the join message
					DeleteMessage(message)
					// Kick (but not ban) the user
					KickUser(message.Chat, user, false)
					continue
				}

				if cfg.KickOnFullNameLength && fullNameLength > cfg.FullNameMaxLength {
					// delete the join message
					DeleteMessage(message)
					// Kick (but not ban) the user
					KickUser(message.Chat, user, false)
				}
			}
		}

	}
}

func PostHelp(chat *tgbotapi.Chat, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(chat.ID,
		"/help visualizza questo messaggio\n"+
			"/roadmap Visualizza la roadmap IoTeX")
	msg.DisableNotification = true
	msg.ParseMode = "Markdown"
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func PostRoadmap(chat *tgbotapi.Chat, message *tgbotapi.Message) {
	msg := tgbotapi.NewPhotoUpload(chat.ID, "images/roadmap.png")
	msg.Caption = "Roadmap IoTeX"
	msg.ReplyToMessageID = message.MessageID
	msg.DisableNotification = true
	bot.Send(msg)
}

// Kick the user from the chat. Optionally, ban it
func KickUser(chat *tgbotapi.Chat, user tgbotapi.User, ban bool) {
	kickConfig := tgbotapi.KickChatMemberConfig{}
	kickConfig.ChatID = chat.ID
	kickConfig.UserID = user.ID
	_, err := bot.KickChatMember(kickConfig)
	if err != nil {
		log.Printf("Error kicking user %s: %s", user.UserName, err)
	} else {
		log.Printf("[KICK] Kicked user %s: Name length = %d > %d", user.UserName, len(user.FirstName), cfg.FirstNameMaxLength)
	}

	// User is banned by default when kicked
	if ban == true {
		return
	}

	unBanConfig := tgbotapi.ChatMemberConfig{}
	unBanConfig.ChatID = chat.ID
	unBanConfig.UserID = user.ID
	bot.UnbanChatMember(unBanConfig)
}

// Delete a message
func DeleteMessage(message *tgbotapi.Message)  {
	delConfig := tgbotapi.DeleteMessageConfig{}
	delConfig.ChatID = message.Chat.ID
	delConfig.MessageID = message.MessageID
	_, err := bot.DeleteMessage(delConfig)
	if err != nil {
		log.Printf("Error deleting join message for user %s: %s", message.From.UserName, err)
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
func GetApiKey() string {
	usingJson := false;

	log.Print("looking for Telegram Bot API Key in BOTAPIKEY environment variable...")
	API := os.Getenv("BOTAPIKEY")
	if API == "" {
		log.Print("Not found\n")
		log.Print("looking for Telegram Bot API Key in cfg.json file...")
		API = cfg.ApiKey
		usingJson = true;
	}

	if API == "" {
		log.Fatal("Telegram API Key not found!")
	}
	if usingJson {
		log.Printf("Using Api key found in cfg file:")
	} else {
		log.Printf("Using Api key found in Environment variable")
	}

	return API
}
