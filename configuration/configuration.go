package configuration

import (
	"github.com/tkanos/gonfig"
	"log"
	"os"
)

type Configuration struct {
	KickOnFirstNameLength bool
	FirstNameMaxLength    int
	KickOnFullNameLength  bool
	FullNameMaxLength     int
	PostWelcomeMessage    bool
	WelcomeMessage        string
}

var Cfg Configuration

func init() {

	err := gonfig.GetConf("configuration.json", &Cfg)

	if err != nil {
		log.Fatal(err)
	}

}

func GetApiKey() string {
	API := os.Getenv("BOTAPIKEY")
	if API == "" {
		log.Print("BOTAPIKEY environment variable not set!\n")
		log.Fatal("Telegram API Key not found!")
	}

	return API
}
