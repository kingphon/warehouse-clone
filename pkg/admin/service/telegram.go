package service

import (
	"fmt"
	"git.selly.red/Selly-Modules/logger"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/config"
	"git.selly.red/Selly-Server/warehouse/pkg/admin/server/initialize/telegrambot"
)

const (
	teleParseModeMarkdown = "Markdown"
	MsgRunJobSuccess      = "\U0001F525 Job \U0001F525 \n *%s*"
)

type TelegramBotHelper interface {
	SendMessenger(message string) error
	SendMessageHaveRunJob(jobName string) error
}

// telegramBotHelperImplement ...
type telegramBotHelperImplement struct {
}

// Telegram ...
func Telegram() TelegramBotHelper {
	return telegramBotHelperImplement{}
}

func (t telegramBotHelperImplement) SendMessenger(message string) error {
	err := telegrambot.SendMessenger(config.GetENV().Telegram.Token, &telegrambot.SendMessageReqBody{
		ChatID:    config.GetENV().Telegram.WarehouseID,
		Text:      message,
		ParseMode: teleParseModeMarkdown,
	})
	return err
}

func (t telegramBotHelperImplement) SendMessageHaveRunJob(jobName string) error {
	msg := fmt.Sprintf(MsgRunJobSuccess, jobName)
	err := telegrambot.SendMessenger(config.GetENV().Telegram.Token, &telegrambot.SendMessageReqBody{
		ChatID:    config.GetENV().Telegram.WarehouseID,
		Text:      msg,
		ParseMode: teleParseModeMarkdown,
	})
	if err != nil {
		logger.Debug("SendMessageHaveRunJob", logger.LogData{
			"error": err.Error(),
		})
	}
	return err
}
