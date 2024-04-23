package types

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

type MenuKeyboard interface {
	GetMenuKeyboard(ctx context.Context, b *bot.Bot, handler inline.OnSelect) models.ReplyMarkup
}

type HandleMessage func(ctx context.Context, b *bot.Bot, mes *models.Message) error
