package telegram

import (
	"fmt"

	"gopkg.in/telebot.v3"
)

func (service *Service) HandleText(ctx telebot.Context) error {
	username := ctx.Chat().Username
	spreadsheetId, found := service.mapping[username]
	if !found {
		fmt.Printf("User %v is not supported!", username)
		return nil
	}

	time := ctx.Message().Unixtime
	text := ctx.Message().Text
	err := service.updater.ProcessUpdate(spreadsheetId, time, text)
	if err == nil {
		ctx.Reply("Принято")
	}
	return err
}
