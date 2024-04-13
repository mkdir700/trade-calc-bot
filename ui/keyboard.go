package ui

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"

	// "capital_calculator_tgbot/task"
)


func CloseButton(ctx context.Context, b *bot.Bot) *inline.Keyboard {
	return inline.New(b).Button(
		"取消",
		[]byte("close"),
		func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			// 删除任务
			// task.GetTaskManager().RemoveTask(mes.Message.Chat.ID)
		},
	)
}
