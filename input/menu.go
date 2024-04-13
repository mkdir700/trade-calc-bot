package input

import (
	t "capital_calculator_tgbot/task"
	"capital_calculator_tgbot/ui"
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type MainMenu struct{}

func NewMainMenu() *MainMenu {
	return &MainMenu{}
}

func (m *MainMenu) SendMessage(ctx context.Context, b *bot.Bot, chatID int64) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "请选择功能",
		ReplyMarkup: ui.MenuKeyboard(ctx, b, m.HandleMessage),
	})
	return err
}

// 处理用户的选择
func (m *MainMenu) HandleMessage(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	tm := t.GetTaskManager()
	switch string(data) {
	case "1-1":
		task := t.NewTask(mes.Message.Chat.ID)
		task.NextStep()
		tm.AddTask(task)
		err := NewInputCapital().SendMessage(ctx, b, mes.Message.Chat.ID)
		if err != nil {
			fmt.Println(err)
		}
	case "1-2":
		// TODO: 未实现
	}
}
