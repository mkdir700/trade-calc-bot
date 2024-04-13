package input

import (
	"context"
	"strconv"

	t "capital_calculator_tgbot/task"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type InputCapital struct {
	Content string
}

func NewInputCapital() *InputCapital {
	return &InputCapital{
		Content: "请输入本金(单位: U), 例如: 1000",
	}
}

func (m *InputCapital) SendMessage(ctx context.Context, b *bot.Bot, chatID int64) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   m.Content,
	})
	return err
}

func (m *InputCapital) HandleMessage(ctx context.Context, b *bot.Bot, mes *models.Message) error {
	task := ctx.Value("task").(*t.Task)
	capital, err := strconv.ParseFloat(mes.Text, 64)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Chat.ID,
			Text:   "请输入有效的数字",
		})
		return err
	}
	if capital <= 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Chat.ID,
			Text:   "请输入大于0的数字",
		})
		return err
	}
	task.Payload.SetCapital(capital)
	return nil
}
