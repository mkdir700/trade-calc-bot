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
	CmdReturnOpenPositionMenu = "return_open_position_menu"
	CmdReturnAddPositionMenu  = "return_add_position_menu"
	CmdReturnMainMenu         = "return_main_menu"
)

type Back struct {
	text              string
	prefix            string
	callbackHandlerID string
	cmd               string // 返回的命令
}

func NewBack(text, cmd string) *Back {
	// 判断 cmd 是否有效
	if cmd != CmdReturnOpenPositionMenu && cmd != CmdReturnAddPositionMenu && cmd != CmdReturnMainMenu {
		log.Println("[Back] [ERROR] Invalid cmd: ", cmd)
		return nil
	}
	return &Back{
		text:   text,
		prefix: bot.RandomString(8),
		cmd:    cmd,
	}
}

func (self *Back) Prefix() string {
	return self.prefix
}

func (self *Back) onError(err error) {
	log.Println("[Back] [ERROR] Error: ", err)
}

func (self *Back) Show(ctx context.Context, b *bot.Bot, chatID any) (*models.Message, error) {
	self.callbackHandlerID = b.RegisterHandler(bot.HandlerTypeCallbackQueryData, self.prefix, bot.MatchTypePrefix, self.callback)

	return b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        self.text,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: self.buildKeyboard(),
	})
}

func (self *Back) callbackAnswer(ctx context.Context, b *bot.Bot, callbackQuery *models.CallbackQuery) {
	ok, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
	})
	if err != nil {
		self.onError(err)
		return
	}
	if !ok {
		self.onError(fmt.Errorf("callback answer failed"))
	}
}

func (self *Back) callback(ctx context.Context, b *bot.Bot, update *models.Update) {
	cmd := strings.TrimPrefix(update.CallbackQuery.Data, self.prefix)
	switch cmd {
	case CmdReturnOpenPositionMenu:
		self.callbackAnswer(ctx, b, update.CallbackQuery)
		_, err := NewOpenPositionMenu().ReplaceShow(
			ctx,
			b,
			update.CallbackQuery.Message.Message.Chat.ID,
			update.CallbackQuery.Message.Message.ID,
			update.CallbackQuery.InlineMessageID,
		)
		if err != nil {
			self.onError(err)
			return
		}
		b.UnregisterHandler(self.callbackHandlerID)
	}
}

func (self *Back) buildKeyboard() models.ReplyMarkup {
	row1 := []models.InlineKeyboardButton{
		{
			Text:         "返回",
			CallbackData: self.prefix + self.cmd,
		},
	}

	kb := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			row1,
		},
	}
	return kb
}
