package main

import (
	"capital_calculator_tgbot/utils"
	"context"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

func main() {
	var err error
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	systemEnv := make(map[string]string)
	for _, env := range os.Environ() {
		kv := strings.Split(env, "=")
		systemEnv[kv[0]] = kv[1]
	}

	if utils.IsProduction() {
		err = godotenv.Load(".env")
	} else {
		err = godotenv.Load(".env.development")
	}
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	for k, v := range systemEnv {
		os.Setenv(k, v)
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

// 默认处理函数
func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	if update.Message.Text == "/start" {
		if !taskManager.HasTask(update.Message.Chat.ID) {
			t := NewTask(update.Message.Chat.ID)
			taskManager.AddTask(t)
		}
		NewMainMenu().Show(ctx, b, update.Message.Chat.ID)
		return
	}
}
