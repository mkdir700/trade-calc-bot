package input

import (
	t "capital_calculator_tgbot/task"
	"context"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type InputLossRatio struct {
	Content string
}

func NewInputLossRatio() *InputLossRatio {
	return &InputLossRatio{
		Content: "请输入亏损比例(单位: %), 例如: 0.38",
	}
}

func (m *InputLossRatio) SendMessage(ctx context.Context, b *bot.Bot, chatId int64) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   m.Content,
	})
	return err
}

func (m *InputLossRatio) HandleMessage(ctx context.Context, b *bot.Bot, mes *models.Message) error {
	task := t.GetTaskManager().GetTask(mes.Chat.ID)
	lossRatio, err := strconv.ParseFloat(string(mes.Text), 64)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Chat.ID,
			Text:   "请输入有效的数字",
		})
		return err
	}
	if lossRatio <= 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Chat.ID,
			Text:   "请输入大于0的数字",
		})
		return err
	}
	if lossRatio > 1 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Chat.ID,
			Text:   "请输入小于1的数字",
		})
		return err
	}
	task.Payload.SetLossRatio(lossRatio / 100)
	return nil
}
