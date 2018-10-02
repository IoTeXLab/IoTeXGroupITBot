package spamfilter

import (
	"time"
	"github.com/IoTeXGroupIT/IoTeXGroupITBot/botApi"
	"github.com/IoTeXGroupIT/IoTeXGroupITBot/configuration"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// EnabledJoinFilter Enable or disable newly joined users spam filter
var EnabledJoinFilter = true;
// EnabledMediaFilter Enable or disable deleting messages containing media or links
var EnabledMediaFilter = false;

var whiteList = []string{"zimne"}
var joinedUsers = make(map[string]time.Time)

// FilterNewUserJoined ban newly joined users if their name/full name contain too much text
func FilterNewUserJoined(bot *tgbotapi.BotAPI, joinMessage *tgbotapi.Message) bool {
	
	user := joinMessage.From
    joinedTime := time.Now()
	joinedUsers[user.UserName] = joinedTime

	if (EnabledJoinFilter == false) {
		return true
	}

	// Spam filter rule applied to First Name length
	firstNameLength := len(user.FirstName)
	// TODO: if we want sto stay even safer, whe could kick only chinese text?
	if configuration.Cfg.KickOnFirstNameLength && firstNameLength > configuration.Cfg.FirstNameMaxLength {
		// delete the join message
		botApi.DeleteMessage(bot, joinMessage)
		// Kick (but not ban) the user
		botApi.KickUser(bot, joinMessage.Chat, user, false)
		return false
	}

	// Spam filter rule applied on Full Name length
	fullNameLength := firstNameLength + len(user.LastName)
	if configuration.Cfg.KickOnFullNameLength && fullNameLength > configuration.Cfg.FullNameMaxLength {
		// delete the join message
		botApi.DeleteMessage(bot, joinMessage)
		// Kick (but not ban) the user
		botApi.KickUser(bot, joinMessage.Chat, user, false)
		return false
	}

	return true
}

// FilterMessageWithLinks delete messages if they contains links
// Works both for direct messages and forwarded messages
func FilterMessageWithLinks(bot *tgbotapi.BotAPI, message *tgbotapi.Message) bool {
	if (EnabledMediaFilter == false) {
		return true;
	}
	
	if (isWhitelisted(message.From)) {
		return true
	}

	if (message.Entities != nil && CanUserPostMedia(message.From)) {
		return true
	}
	
	joinedTime := joinedUsers[message.From.UserName]
	elapsedMinutes := time.Since(joinedTime).Minutes()
	
	if (elapsedMinutes < 60) {
		// delete the message
		botApi.DeleteMessage(bot, message)
		// Kick (but not ban) the user
		botApi.KickUser(bot, message.Chat, message.From, false)
		return false
	}

	return true;
}

// CanUserPostMedia returns true if the user is allowed to post media files/links
func CanUserPostMedia(user *tgbotapi.User) bool {
	/* 
		Implement here the rules for a user to be allowed to post media/links
	*/
	return false
}	

func isWhitelisted(user *tgbotapi.User) bool {
	for i:=0;i<len(whiteList);i++ {
		if (whiteList[i] == user.UserName) {
			return true
		}
	}
	return false
}

