package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	cmd0_01   = "0.01"
	cmd0_02   = "0.02"
	cmd0_03   = "0.03"
	cmdCustom = "custom"
)

type InputCapitalLossRadio struct {
	text              string
	prefix            string
	callbackHandlerID string
}

func NewInputCapitalLossRadio() *InputCapitalLossRadio {
	return &InputCapitalLossRadio{
		text:   "请选择本金亏损比例",
		prefix: bot.RandomString(8),
	}
}

func (m *InputCapitalLossRadio) Prefix() string {
	return m.prefix
}

func (m *InputCapitalLossRadio) buildText() string {
	return m.text
}

func (m *InputCapitalLossRadio) buildKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "1%",
					CallbackData: m.prefix + "0.01",
				},
				{
					Text:         "2%",
					CallbackData: m.prefix + "0.02",
				},
				{
					Text:         "3%",
					CallbackData: m.prefix + "0.03",
				},
			},
			// {
			// 	{
			// 		Text:         "自定义",
			// 		CallbackData: m.prefix + "custom",
			// 	},
			// },
		},
	}
}

func (m *InputCapitalLossRadio) Show(ctx context.Context, b *bot.Bot, chatID int64) error {
	m.callbackHandlerID = b.RegisterHandler(bot.HandlerTypeCallbackQueryData, m.prefix, bot.MatchTypePrefix, m.callback)
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        m.buildText(),
		ReplyMarkup: m.buildKeyboard(),
	})
	return err
}

func (m *InputCapitalLossRadio) onError(err error) {
	log.Printf("[InputCapitalLossRadio] [ERROR] Error: %v", err)
}

func (m *InputCapitalLossRadio) callbackAnswer(ctx context.Context, b *bot.Bot, callbackQuery *models.CallbackQuery) {
	ok, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
	})
	if err != nil {
		m.onError(err)
		return
	}
	if !ok {
		m.onError(fmt.Errorf("callback answer failed"))
	}
}

func (m *InputCapitalLossRadio) callback(ctx context.Context, b *bot.Bot, update *models.Update) {
	cmd := strings.TrimPrefix(update.CallbackQuery.Data, m.prefix)
	task := taskManager.GetTask(update.CallbackQuery.Message.Message.Chat.ID)

	back := func() {
		_, err := NewOpenPositionMenu().ReplaceShow(
			ctx,
			b,
			update.CallbackQuery.Message.Message.Chat.ID,
			update.CallbackQuery.Message.Message.ID,
			update.CallbackQuery.InlineMessageID,
		)
		if err != nil {
			m.onError(err)
			return
		}
	}

	switch cmd {
	case cmd0_01:
		task.Payload.SetCapitalLossRatio(0.01)
		m.callbackAnswer(ctx, b, update.CallbackQuery)
		back()
	case cmd0_02:
		m.callbackAnswer(ctx, b, update.CallbackQuery)
		task.Payload.SetCapitalLossRatio(0.02)
		back()
	case cmd0_03:
		m.callbackAnswer(ctx, b, update.CallbackQuery)
		task.Payload.SetCapitalLossRatio(0.03)
		back()
	case cmdCustom:
		m.onError(fmt.Errorf("custom"))
	}
	b.UnregisterHandler(m.callbackHandlerID)
}

func (m *InputCapitalLossRadio) handleCustomOption(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Message.Chat.ID,
		Text:   "请输入本金损失比例(单位: %), 例如: 2",
	})
}

func (m *InputCapitalLossRadio) HandleMessage(ctx context.Context, b *bot.Bot, mes *models.Message) error {
	capitalLossRatio, err := strconv.ParseFloat(string(mes.Text), 64)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Chat.ID,
			Text:   "请输入有效的数字",
		})
		return err
	}
	if capitalLossRatio <= 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Chat.ID,
			Text:   "请输入大于0的数字",
		})
		return err
	}
	if capitalLossRatio > 100 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Chat.ID,
			Text:   "请输入小于100的数字",
		})
		return err
	}
	task := GetTaskManager().GetTask(mes.Chat.ID)
	task.Payload.SetCapitalLossRatio(capitalLossRatio / 100)
	return nil
}
