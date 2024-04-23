package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type OnErrorHandler func(err error)

const (
	cmdInputCapital          = "input_capital"
	cmdInputCapitalLossRatio = "input_capital_loss_ratio"
	cmdInputLossRatio        = "input_loss_ratio"
	cmdCalculate             = "calculate"
	cmdCancel                = "cancel"
)

type OpenPositionMenu struct {
	closeButton       string
	onError           OnErrorHandler
	prefix            string
	messageID         int
	callbackHandlerID string
}

func NewOpenPositionMenu() *OpenPositionMenu {
	return &OpenPositionMenu{
		closeButton: "取消",
		onError:     defaultOnError,
		prefix:      bot.RandomString(8),
	}
}

func (o *OpenPositionMenu) Prefix() string {
	return o.prefix
}

func defaultOnError(err error) {
	log.Printf("[OpenPosition] [ERROR] Error: %v", err)
}

func (o *OpenPositionMenu) Show(ctx context.Context, b *bot.Bot, chatID int64) (*models.Message, error) {
	o.callbackHandlerID = b.RegisterHandler(bot.HandlerTypeCallbackQueryData, o.prefix, bot.MatchTypePrefix, o.callback)

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        o.buildText(chatID),
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: o.buildKeyboard(),
	})
	o.messageID = msg.ID
	return msg, err
}

func (o *OpenPositionMenu) ReplaceShow(ctx context.Context, b *bot.Bot, chatID int64, messageID int, inlineMessageID string) (*models.Message, error) {
	o.callbackHandlerID = b.RegisterHandler(bot.HandlerTypeCallbackQueryData, o.prefix, bot.MatchTypePrefix, o.callback)
	o.messageID = messageID

	return b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:          chatID,
		MessageID:       messageID,
		InlineMessageID: inlineMessageID,
		Text:            o.buildText(chatID),
		ParseMode:       models.ParseModeHTML,
		ReplyMarkup:     o.buildKeyboard(),
	})
}

func (o *OpenPositionMenu) buildKeyboard() models.InlineKeyboardMarkup {
	row1 := []models.InlineKeyboardButton{
		{
			Text:         "本金",
			CallbackData: o.prefix + cmdInputCapital,
		},
	}

	// 第二行：输入本金亏损比例按钮
	row2 := []models.InlineKeyboardButton{
		{
			Text:         "本金亏损比例",
			CallbackData: o.prefix + cmdInputCapitalLossRatio,
		},
	}

	// 第三行：输入亏损比例按钮
	row3 := []models.InlineKeyboardButton{
		{
			Text:         "本单亏损比例",
			CallbackData: o.prefix + cmdInputLossRatio,
		},
	}

	// 计算按钮
	row4 := []models.InlineKeyboardButton{
		{
			Text:         "计算",
			CallbackData: o.prefix + cmdCalculate,
		},
		{
			Text:         "取消",
			CallbackData: o.prefix + cmdCancel,
		},
	}

	kb := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			row1,
			row2,
			row3,
			row4,
		},
	}
	return kb
}

func (o *OpenPositionMenu) callbackAnswer(ctx context.Context, b *bot.Bot, callbackQuery *models.CallbackQuery, text string) {
	ok, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
		Text:            text,
	})
	if err != nil {
		o.onError(err)
		return
	}
	if !ok {
		o.onError(fmt.Errorf("callback answer failed"))
	}
}

// 删除菜单
func (o *OpenPositionMenu) deleteSelf(ctx context.Context, b *bot.Bot, update *models.Update) (bool, error) {
	return b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: o.messageID,
	})
}

func (o *OpenPositionMenu) callback(ctx context.Context, b *bot.Bot, update *models.Update) {
	cmd := strings.TrimPrefix(update.CallbackQuery.Data, o.prefix)
	switch cmd {
	case cmdInputCapital:
		o.deleteSelf(ctx, b, update)
		NewInputCapital().Show(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		o.callbackAnswer(ctx, b, update.CallbackQuery, "")
	case cmdInputCapitalLossRatio:
		o.deleteSelf(ctx, b, update)
		NewInputCapitalLossRadio().Show(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		o.callbackAnswer(ctx, b, update.CallbackQuery, "")
	case cmdInputLossRatio:
		o.deleteSelf(ctx, b, update)
		NewInputLossRatio().Show(ctx, b, update.CallbackQuery.Message.Message.Chat.ID)
		o.callbackAnswer(ctx, b, update.CallbackQuery, "")
	case cmdCalculate:
		task := GetTaskManager().GetTask(update.CallbackQuery.Message.Message.Chat.ID)
		if task == nil {
			o.onError(errors.New("无法获取任务"))
			o.callbackAnswer(ctx, b, update.CallbackQuery, "无法获取任务")
			return
		}
		if task.Payload == nil {
			o.onError(errors.New("无法获取任务数据"))
			o.callbackAnswer(ctx, b, update.CallbackQuery, "无法获取任务数据")
			return
		}
		result := NewOpenPositionResult(*task.Payload)
		err := result.Calculate()
		if err != nil {
			o.onError(err)
			o.callbackAnswer(ctx, b, update.CallbackQuery, err.Error())
			return
		}
		text := result.BuildText()
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   text,
		})
		o.callbackAnswer(ctx, b, update.CallbackQuery, "")
		o.deleteSelf(ctx, b, update)
		b.UnregisterHandler(o.callbackHandlerID)
	case cmdCancel:
		o.callbackAnswer(ctx, b, update.CallbackQuery, "")
		GetTaskManager().RemoveTask(update.CallbackQuery.Message.Message.Chat.ID)
		_, errDelete := o.deleteSelf(ctx, b, update)
		if errDelete != nil {
			o.onError(errDelete)
		}
		b.UnregisterHandler(o.callbackHandlerID)
	}
}

func (o *OpenPositionMenu) buildText(chatId int64) string {
	task := GetTaskManager().GetTask(chatId)
	if task == nil {
		return "无法获取任务"
	}
	capital := task.Payload.Capital
	capitalLossRatio := task.Payload.CapitalLossRatio
	lossRatio := task.Payload.LossRatio

	unsetText := "未设置"
	textCapital := unsetText
	if capital != 0 {
		textCapital = fmt.Sprintf("%0.2f", capital)
	}

	textCapitalLossRatio := unsetText
	if capitalLossRatio != 0 {
		// 0.02 -> 2%
		textCapitalLossRatio = fmt.Sprintf("%0.2f%%", capitalLossRatio*100)
	}

	textLossRatio := unsetText
	if lossRatio != 0 {
		textLossRatio = fmt.Sprintf("%0.2f%%", lossRatio*100)
	}

	template := `计算开仓

本金: %s

本金亏损比例: %s

本单亏损比例: %s

点击下方按钮输入相关数据`
	text := fmt.Sprintf(
		template,
		textCapital,
		textCapitalLossRatio,
		textLossRatio,
	)
	fmt.Println(text)
	return text
}
