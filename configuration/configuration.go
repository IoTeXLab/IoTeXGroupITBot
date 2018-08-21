package configuration

import (
	"github.com/tkanos/gonfig"
	"log"
	"os"
)

type Configuration struct {
	ApiKey                string
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
	usingJson := false;
	log.Print("looking for Telegram Bot API Key in BOTAPIKEY environment variable...")
	API := os.Getenv("BOTAPIKEY")
	if API == "" {
		log.Print("Not found\n")
		log.Print("looking for Telegram Bot API Key in cfg.json file...")
		API = Cfg.ApiKey
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
