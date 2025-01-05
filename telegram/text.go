package telegram

import (
	"gopkg.in/telebot.v3"
)

func (service *Service) HandleText(ctx telebot.Context) error {
	username := ctx.Chat().Username
	spreadsheetId, found := service.mapping[username]
	if !found {
		logger.Infof("User %v is not supported!", username)
		return ctx.Reply("Неизвестный пользователь")
	}

	time := ctx.Message().Unixtime
	text := ctx.Message().Text
	err := service.updater.ProcessUpdate(spreadsheetId, time, text)
	if err == nil {
		ctx.Reply("Принято")
	}
	return err
}
