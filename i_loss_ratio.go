package main

import (
	"capital_calculator_tgbot/utils"
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type InputLossRatio struct {
	text              string
	callbackHandlerID string
}

func NewInputLossRatio() *InputLossRatio {
	return &InputLossRatio{
		text: "请输入止损比例(单位: %), 例如: 0.38",
	}
}

func (m *InputLossRatio) SendMessage(ctx context.Context, b *bot.Bot, chatId int64) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   m.text,
	})
	return err
}

func (m *InputLossRatio) Show(ctx context.Context, b *bot.Bot, chatID int64) error {
	// 注册消息回调函数
	m.callbackHandlerID = b.RegisterHandlerMatchFunc(
		func(update *models.Update) bool {
			return true
		},
		m.callback,
	)
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   m.text,
		ReplyMarkup: &models.ForceReply{
			ForceReply: true,
		},
	})
	return err
}

func (m *InputLossRatio) callback(ctx context.Context, b *bot.Bot, update *models.Update) {
	task := GetTaskManager().GetTask(update.Message.Chat.ID)
	if task == nil {
		return
	}
	lossRatio, err := utils.ParseFloat(update.Message.Text)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "请输入有效的数字",
		})
		return
	}

	task.Payload.SetLossRatio(lossRatio / 100)

	// 删除用户的回复消息
	_, err = b.DeleteMessages(ctx, &bot.DeleteMessagesParams{
		ChatID:     update.Message.Chat.ID,
		MessageIDs: []int{update.Message.ID, update.Message.ReplyToMessage.ID},
	})
	
	if err != nil {
		m.onError(err)
		return
	}

	// 返回上一级菜单
	_, err = NewOpenPositionMenu().Show(
		ctx,
		b,
		update.Message.Chat.ID,
	)

	if err != nil {
		m.onError(err)
		return
	}
	b.UnregisterHandler(m.callbackHandlerID)
}

func (m *InputLossRatio) onError(err error) {
	log.Println("[InputLossRatio] [ERROR] Error: ", err)
}