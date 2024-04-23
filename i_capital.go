package main

import (
	"capital_calculator_tgbot/utils"
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type InputCapital struct {
	text              string
	prefix            string
	callbackHandlerID string
	Callback          func(ctx context.Context, b *bot.Bot, update *models.Update)
}

func NewInputCapital() *InputCapital {
	return &InputCapital{
		text: "请输入本金(单位: U), 例如: 1000",
	}
}

func (m *InputCapital) Prefix() string {
	return m.prefix
}

func (m *InputCapital) BuildText(text string) {
	m.text = text
}

func (m *InputCapital) matchFunc(update *models.Update) bool {
	// fmt.Println(update.Message.ReplyToMessage.ID)
	// if update.Message.ReplyToMessage == nil {
	// 	return false
	// }
	// return m.botMessageID == update.Message.ReplyToMessage.ID
	return true
}

func (m *InputCapital) Show(ctx context.Context, b *bot.Bot, chatID int64) error {
	// 注册消息回调函数
	m.callbackHandlerID = b.RegisterHandlerMatchFunc(m.matchFunc, m.callback)
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   m.text,
		// ReplyMarkup: &models.ForceReply{
		// 	ForceReply: true,
		// },
	})
	return err
}

func (m *InputCapital) callback(ctx context.Context, b *bot.Bot, update *models.Update) {
	taskManager := GetTaskManager()
	task := taskManager.GetTask(update.Message.Chat.ID)
	val, err := utils.ParseFloat(update.Message.Text)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "请输入有效的数字",
		})
		return
	}
	if err = task.Payload.SetCapital(val); err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   err.Error(),
		})
		return
	}

	NewBack(
		"本金已设置完成，请返回上一级菜单",
		CmdReturnOpenPositionMenu,
	).Show(ctx, b, update.Message.Chat.ID)

	b.UnregisterHandler(m.callbackHandlerID)
}
