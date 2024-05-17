package main

import (
	"log"
	"capital_calculator_tgbot/utils"
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type InputCapital struct {
	text              string
	prefix            string
	callbackHandlerID string
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

func (m *InputCapital) onError(err error) {
	log.Println("[InputCapital] [ERROR] Error: ", err)
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
	m.callbackHandlerID = b.RegisterHandlerMatchFunc(m.matchFunc, m.callback)
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   m.text,
		ReplyMarkup: &models.ForceReply{
			ForceReply: true,
		},
	})
	return err
}

func (m *InputCapital) ReplaceShow(ctx context.Context, b *bot.Bot, chatID int64, messageID int, inlineMessageID string) (*models.Message, error) {
	m.callbackHandlerID = b.RegisterHandler(bot.HandlerTypeCallbackQueryData, m.prefix, bot.MatchTypePrefix, m.callback)
	return b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:          chatID,
		MessageID:       messageID,
		InlineMessageID: inlineMessageID,
		Text:            m.text,
		ParseMode:       models.ParseModeHTML,
		ReplyMarkup: &models.ForceReply{
			ForceReply: true,
		},
	})
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
