package main

import (
	"capital_calculator_tgbot/input"
	t "capital_calculator_tgbot/task"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

func main() {
	var err error
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if IsProduction() {
		err = godotenv.Load()
	} else {
		err = godotenv.Load(".env.development")
	}
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("BOT_TOKEN")
	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var err error
	taskManager := t.GetTaskManager()
	if !taskManager.HasTask(update.Message.Chat.ID) {
		// 无任务
		err = input.NewMainMenu().SendMessage(ctx, b, update.Message.Chat.ID)
		if err != nil {
			fmt.Println(err)
		}
		return
	} else {
		// 有任务
		handleMessage(ctx, b, update.Message)
	}
}

func handleMessage(
	ctx context.Context,
	b *bot.Bot,
	mes *models.Message,
) error {
	var err error
	task := t.GetTaskManager().GetTask(mes.Chat.ID)
	// 将当前任务绑定在上下文中
	ctx = context.WithValue(ctx, "task", task)
	// TODO: 应该采用什么设计模式来执行不同任务的步骤？
	switch task.Step {
	// case t.Start:
	// 	// input.NewMainMenu().HandleMessage(ctx, b, mes)
	// 	input.NewInputCapital().SendMessage(ctx, b, mes.Chat.ID)
	// 	task.NextStep()
	case t.InputCapital:
		err = input.NewInputCapital().HandleMessage(ctx, b, mes)
		if err != nil {
			return err
		}
		err = input.NewInputCapitalLossRadio().SendMessage(ctx, b, mes.Chat.ID)
		if err != nil {
			return err
		}
		if task.Step == t.InputCapital {
			task.NextStep()
		}
	case t.InputCapitalLossRadio:
		err = input.NewInputCapitalLossRadio().HandleMessage(ctx, b, mes)
		if err != nil {
			return err
		}
		err = input.NewInputLossRatio().SendMessage(ctx, b, mes.Chat.ID)
		if err != nil {
			return err
		}
		task.NextStep()
	case t.InputLossRadio:
		err = input.NewInputLossRatio().HandleMessage(ctx, b, mes)
		if err != nil {
			return err
		}
		task.NextStep()
	}

	if t.End == task.Step {
		// 计算
		result := t.NewOpenPositionResult(*task.Payload)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Chat.ID,
			Text:   result.ShowMessage(),
		})
		t.GetTaskManager().RemoveTask(mes.Chat.ID)
	}
	return nil
}
