package reminder

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
	"github.com/simonerom/IoTeXGroupITBot/botApi"
	"strconv"
	"sort"
)

type ReminderClass int

const (
	// This event has only a start date/time, it happens at a specific moment (e.g. meetings)
	OneShot ReminderClass = 0
	// This event has a start and finish date/time, it happens during an interval of time (e.g. bounty programs)
	StartEndDate ReminderClass = 1
	// This event has a only a finish date/time it ends by a specific moment (e.g. deadlines)
	Deadline ReminderClass = 2
)

type Reminder struct {
	title       string
	startDate   time.Time
	finishDate  time.Time
	reminderMsg string
	startMsg    string
	finishMsg   string
}

var reminders []Reminder
var PST *time.Location

func init() {
	PST, _ = time.LoadLocation("America/Mazatlan")
	var reminder1 = new(Reminder)
	var reminder2 = new(Reminder)
	var reminder3 = new(Reminder)
	reminder1.title = "Rilascio Testnet Beta"
	reminder1.finishDate = time.Date(2018, time.August, 31, 23, 59, 59, 0, PST)
	reminder2.title = "Sblocco Tokens #2"
	reminder2.finishDate = time.Date(2018, time.August, 22, 23, 59, 59, 0, PST)
	reminder3.title = "Sessione AMA"
	reminder3.finishDate = time.Date(2018, time.August, 18, 17, 0, 0, 0, PST)

	reminders = append(append(append(reminders, *reminder1), *reminder2), *reminder3)
	sort.Slice(reminders, func(i, j int) bool {
		return reminders[i].finishDate.Unix() < reminders[j].finishDate.Unix()
	})
}

// Posts a list of the next events
func List(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat) {
	msg := "*Prossimi Eventi IoTeX*"

	var previousLeftTime time.Duration
	for _, reminder := range (reminders) {
		leftTime, leftTimeStr := getLeftTime(reminder.finishDate)

		if leftTime < 0 {
			continue
		}

		if (leftTime != previousLeftTime) {
			msg += "\n\n*tra " + leftTimeStr + "*"
		}
		msg += "\n" + reminder.title
		previousLeftTime = leftTime
	}
	botApi.PostTextMessage(bot, chat, msg, false, nil)
}

func getLeftTime(date time.Time) (time.Duration, string) {
	left := date.Sub(time.Now())

	hours := left.Hours()

	if hours > 24 {
		return left, strconv.FormatFloat(hours/24, 'f', 0, 64) + " giorni"
	}

	mins := left.Minutes()

	if mins > 60 {
		return left, strconv.FormatFloat(mins/60, 'f', 0, 64) + " ore " + strconv.FormatInt(int64(mins)%int64(60), 10) + " min "
	}

	return left, strconv.FormatFloat(mins, 'f', 1, 64) + " minuti"
}
