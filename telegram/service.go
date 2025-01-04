package telegram

import (
	"fmn/journalbot/updater"
	"log"
	"os"
	"time"

	"gopkg.in/telebot.v3"
)

type Service struct {
	bot     *telebot.Bot
	mapping map[string]string
	updater *updater.Updater
}

func InitService(updater *updater.Updater) *Service {
	service := Service{
		bot:     initBot(),
		mapping: loadSheetsMapping(),
		updater: updater,
	}
	service.initHandlers()
	return &service
}

func initBot() *telebot.Bot {
	token := os.Getenv("BOT_TOKEN")
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return b
}

func loadSheetsMapping() map[string]string {
	result := make(map[string]string)
	result[os.Getenv("FIRST_USER_USERNAME")] = os.Getenv("FIRST_USER_SPREADSHEET_ID")
	return result
}

func (service *Service) initHandlers() {
	service.bot.Handle(telebot.OnText, service.HandleText)
}

func (service *Service) Start() {
	service.bot.Start()
}
