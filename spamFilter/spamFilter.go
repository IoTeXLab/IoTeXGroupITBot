package spamFilter

import (
	"github.com/IoTeXGroupIT/IoTeXGroupITBot/botApi"
	"github.com/IoTeXGroupIT/IoTeXGroupITBot/configuration"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func FilterNewUserJoined(bot *tgbotapi.BotAPI, joinMessage *tgbotapi.Message) bool {
	user := joinMessage.From

	// Spam filter rule applied to First Name length
	firstNameLength := len(user.FirstName)
	// TODO: if we want sto stay even safer, whe could kick only chinese text?
	if configuration.Cfg.KickOnFirstNameLength && firstNameLength > configuration.Cfg.FirstNameMaxLength {
		// delete the join message
		botApi.DeleteMessage(bot, joinMessage)
		// Kick (but not ban) the user
		botApi.KickUser(bot, joinMessage.Chat, *user, false)
		return false
	}

	// Spam filter rule applied on Full Name length
	fullNameLength := firstNameLength + len(user.LastName)
	if configuration.Cfg.KickOnFullNameLength && fullNameLength > configuration.Cfg.FullNameMaxLength {
		// delete the join message
		botApi.DeleteMessage(bot, joinMessage)
		// Kick (but not ban) the user
		botApi.KickUser(bot, joinMessage.Chat, *user, false)
		return false
	}

	return true
}
