package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	cmdOpenPosition = "1-1"
	cmdAddPosition  = "1-2"
)

type MainMenu struct {
	openPositionMenu  *OpenPositionMenu
	prefix            string
	callbackHandlerID string
}

func NewMainMenu() *MainMenu {
	return &MainMenu{
		openPositionMenu: NewOpenPositionMenu(),
		prefix:           bot.RandomString(8),
	}
}

func (m *MainMenu) Prefix() string {
	return m.prefix
}

func (m *MainMenu) OnError(err error) {
	log.Printf("[MainMenu] [ERROR] Error: %v", err)
}

func (m *MainMenu) buildText() string {
	return `*主菜单*

点击下方按钮开始使用

`
}

func (m *MainMenu) Show(ctx context.Context, b *bot.Bot, chatID any) (*models.Message, error) {
	m.callbackHandlerID = b.RegisterHandler(bot.HandlerTypeCallbackQueryData, m.prefix, bot.MatchTypePrefix, m.callback)

	return b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        m.buildText(),
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: m.buildKeyboard(),
	})
}

func (m *MainMenu) onError(err error) {
	log.Printf("[InputCapitalLossRadio] [ERROR] Error: %v", err)
}

func (m *MainMenu) callbackAnswer(ctx context.Context, b *bot.Bot, callbackQuery *models.CallbackQuery) {
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

func (m *MainMenu) callback(ctx context.Context, b *bot.Bot, update *models.Update) {
	cmd := strings.TrimPrefix(update.CallbackQuery.Data, m.prefix)
	switch cmd {
	case cmdOpenPosition:
		m.callbackAnswer(ctx, b, update.CallbackQuery)
		_, err := m.openPositionMenu.ReplaceShow(
			ctx,
			b,
			update.CallbackQuery.Message.Message.Chat.ID,
			update.CallbackQuery.Message.Message.ID,
			update.CallbackQuery.InlineMessageID,
		)
		if err != nil {
			m.OnError(err)
		}
	case cmdAddPosition:
		// TODO
	}
	b.UnregisterHandler(m.callbackHandlerID)
}

func (m *MainMenu) buildKeyboard() models.ReplyMarkup {
	row1 := []models.InlineKeyboardButton{
		{
			Text:         "计算开仓",
			CallbackData: m.prefix + cmdOpenPosition,
		},
	}

	row2 := []models.InlineKeyboardButton{
		{
			Text:         "计算加仓",
			CallbackData: m.prefix + cmdAddPosition,
		},
	}

	kb := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			row1,
			row2,
		},
	}
	return kb
}
