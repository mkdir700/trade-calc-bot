package input

import (
	t "capital_calculator_tgbot/task"
	"context"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

type InputCapitalLossRadio struct {
	Content string
}

func NewInputCapitalLossRadio() *InputCapitalLossRadio {
	return &InputCapitalLossRadio{
		Content: "请选择本金损失比例",
	}
}

func (m *InputCapitalLossRadio) SendMessage(ctx context.Context, b *bot.Bot, chatId int64) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   m.Content,
		ReplyMarkup: inline.New(b).
			Row().
			Button("1%", []byte("0.01"), m.handleOption).
			Button("2%", []byte("0.02"), m.handleOption).
			Button("3%", []byte("0.03"), m.handleOption).
			Row().
			Button("自定义", []byte("custom"), m.handleCustomOption),
	})
	return err
}

func (m *InputCapitalLossRadio) handleOption(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	val, _ := strconv.ParseFloat(string(data), 64)
	task := t.GetTaskManager().GetTask(mes.Message.Chat.ID)
	task.Payload.SetCapitalLossRatio(val)
	// TODO: 需要解藕，不应该直接调用下一个步骤
	// 考虑针对不同的任务，有不同的调用顺序
	task.NextStep()
	NewInputLossRatio().SendMessage(ctx, b, mes.Message.Chat.ID)
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
	task := t.GetTaskManager().GetTask(mes.Chat.ID)
	task.Payload.SetCapitalLossRatio(capitalLossRatio / 100)
	return nil
}
