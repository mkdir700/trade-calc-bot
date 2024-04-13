package ui

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

func MenuKeyboard(ctx context.Context, b *bot.Bot, handler inline.OnSelect) models.ReplyMarkup {
	return inline.New(b).
		Row().
		Button("开仓计算", []byte("1-1"), handler).
		Button("加仓计算", []byte("1-2"), handler)
}
