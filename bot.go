package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"github.com/tkanos/gonfig"
	"time"
)

type Configuration struct {
	ApiKey string
	KickOnFirstNameLength bool
	FirstNameMaxLength int
	KickOnFullNameLength bool
	FullNameMaxLength int
	CommandsWhitelist []string
}

var configuration Configuration

func main() {

	err := gonfig.GetConf("configuration.json", &configuration)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(configuration)

	API := GetApiKey();

	log.Printf("Found Telegram Bot API: %s", API)
	bot, err := tgbotapi.NewBotAPI(API)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	updates.Clear()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			// Process commands here
		}

		// Filter spam by kicking new users when they include long text as First/Last name
		if update.Message.NewChatMembers != nil {
			for _,user:=range *update.Message.NewChatMembers {
				log.Printf("_________________________")
				log.Printf(time.Now().Format("01-02-2018 15:04:05"))
				log.Printf("_________________________")
				log.Printf("New user joined the group")
				log.Printf("Username: %s", user.UserName)
				log.Printf("First Name: %s",user.FirstName)
				log.Printf("Last Name: %s", user.LastName)
				log.Printf("_________________________")

				if configuration.KickOnFirstNameLength && len(user.FirstName) > configuration.FirstNameMaxLength {
					// delete the join message
					delConfig := tgbotapi.DeleteMessageConfig{}
					delConfig.ChatID = update.Message.Chat.ID
					delConfig.MessageID = update.Message.MessageID
					_,err := bot.DeleteMessage(delConfig)
					if err != nil{
						log.Printf("Error deleting join message for user %s: %s",user.UserName, err)
					}

					kickConfig:=tgbotapi.KickChatMemberConfig{}
					kickConfig.ChatID = update.Message.Chat.ID
					kickConfig.UserID = user.ID
					_, err = bot.KickChatMember(kickConfig)
					if err != nil{
						log.Printf("Error kicking user %s: %s",user.UserName, err)
					} else {
						log.Printf("[KICK] Kicked user %s: Name length = %d > %d", user.UserName, len(user.FirstName), configuration.FirstNameMaxLength)
					}
				}

				fullNameLength := len(user.FirstName) + len(user.LastName)

				if configuration.KickOnFullNameLength && fullNameLength > configuration.FullNameMaxLength {
					// delete the join message
					delConfig := tgbotapi.DeleteMessageConfig{}
					delConfig.ChatID = update.Message.Chat.ID
					delConfig.MessageID = update.Message.MessageID
					_,err := bot.DeleteMessage(delConfig)
					if err != nil{
						log.Printf("Error deleting join message for user %s: %s",user.UserName, err)
					}

					// kick the user
					kickConfig:=tgbotapi.KickChatMemberConfig{}
					kickConfig.ChatID = update.Message.Chat.ID
					kickConfig.UserID = user.ID
					_, err = bot.KickChatMember(kickConfig)
					if err != nil{
						log.Printf("Error kicking user %s: %s",user.UserName, err)
					} else {
						log.Printf("[KICK] Kicked user %s: Full name length = %d > %d", user.UserName, fullNameLength, configuration.FullNameMaxLength)
					}
				}
			}
		}

	}
}

func GetApiKey() string {
	// "673613669:AAFt-sbY8CA67oRpUPCV5O4P5cYQPFhLzM0"

	usingJson := false;

	log.Print("looking for Telegram Bot API Key in BOTAPIKEY environment variable...")
	API := os.Getenv("BOTAPIKEY")
	if API == "" {
		log.Print("Not found\n")
		log.Print("looking for Telegram Bot API Key in configuration.json file...")
		API = configuration.ApiKey
		usingJson = true;
	}

	if API == "" {
		log.Fatal("Telegram API Key not found!")
	}
	if usingJson {
		log.Printf("Using Api key found in configuration file:")
	} else {
		log.Printf("Using Api key found in Environment variable: ")
	}

	log.Println(API)

	return API
}
