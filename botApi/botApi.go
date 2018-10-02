package botApi

import (
	"fmt"
	"github.com/IoTeXGroupIT/IoTeXGroupITBot/configuration"
	. "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

var lastWelcomeMessage *Message

func DeleteMessage(bot *BotAPI, message *Message) {
	delConfig := DeleteMessageConfig{}
	delConfig.ChatID = message.Chat.ID
	delConfig.MessageID = message.MessageID
	_, err := bot.DeleteMessage(delConfig)
	if err != nil {
		log.Printf("Error deleting join message for user %s: %s", message.From.UserName, err)
	}
}

func PostTextMessage(bot *BotAPI, chat *Chat, text string, notification bool, replyTo *Message) {
	msg := NewMessage(chat.ID, text)
	msg.DisableNotification = !notification
	msg.ParseMode = "MarkDown"
	if replyTo != nil {
		msg.ReplyToMessageID = replyTo.MessageID
	}
	bot.Send(msg)
}

func PostWelcomeMessage(bot *BotAPI, chat *Chat, joinMessage *Message) {
	msg := NewMessage(chat.ID, fmt.Sprintf(configuration.Cfg.WelcomeMessage, joinMessage.From.FirstName))
	msg.ReplyToMessageID = joinMessage.MessageID
	msg.DisableNotification = true
	if lastWelcomeMessage != nil {
		DeleteMessage(bot, lastWelcomeMessage)
	}
	sentMessage, err := bot.Send(msg)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	lastWelcomeMessage = &sentMessage
}
func PostWrongCommandMessage(bot *BotAPI, chat *Chat, message *Message) {
	msg := NewMessage(chat.ID,
		"Il comando inserito non è valido.\n"+
			"Se cerchi il Bounty Bot di IoTeX "+
			"clicca su @IoTeXBountyBot e inizia"+
			"una conversazione privata con il bot.")
	msg.DisableNotification = true
	msg.ParseMode = "Markdown"
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func PostHelpMessage(bot *BotAPI, chat *Chat, message *Message) {
	msg := NewMessage(chat.ID,
		"/help visualizza questo messaggio\n"+
			"/roadmap Visualizza la roadmap IoTeX")
	msg.DisableNotification = true
	msg.ParseMode = "Markdown"
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func PostRoadmapImage(bot *BotAPI, chat *Chat, message *Message) {
	msg := NewPhotoUpload(chat.ID, "images/roadmap.png")
	msg.Caption = "Roadmap IoTeX"
	msg.ReplyToMessageID = message.MessageID
	msg.DisableNotification = true
	bot.Send(msg)
}

// KickUser Kick the user from the chat. Optionally, ban it
func KickUser(bot *BotAPI, chat *Chat, user *User, ban bool) {
	kickConfig := KickChatMemberConfig{}
	kickConfig.ChatID = chat.ID
	kickConfig.UserID = user.ID
	_, err := bot.KickChatMember(kickConfig)
	if err != nil {
		log.Printf("Error kicking user %s: %s", user.UserName, err)
	} else {
		log.Printf("[KICK] Kicked user %s: Name length = %d > %d", user.UserName, len(user.FirstName), configuration.Cfg.FirstNameMaxLength)
	}

	// User is banned by default when kicked
	if ban == true {
		return
	}

	unBanConfig := ChatMemberConfig{}
	unBanConfig.ChatID = chat.ID
	unBanConfig.UserID = user.ID
	bot.UnbanChatMember(unBanConfig)
}
