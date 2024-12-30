package main

import (
	"context"
	"log"
	"os"

	"github.com/tikhonp/alcs/internal/apps/telegram/bot"
	"github.com/tikhonp/alcs/internal/config"
	"github.com/tikhonp/alcs/internal/db"
	"github.com/tikhonp/alcs/internal/router"
	"github.com/tikhonp/alcs/internal/util/annalist"
)

func main() {

	// Read the configuration from the pkl file
	cfg, err := config.LoadFromPath(context.Background(), "config.pkl")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Connect to the database
	modelsFactory, err := db.Connect(cfg.Db)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Annalist
	annalistManager := annalist.NewDefaultAnnalist(cfg.Server.Debug)

	// Initialize Telegram bot
	bot, err := bot.NewBot(cfg.Auth.Telegram.BotToken)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}
	// TODO: Activate it for "GOOD" domain
	// bot.SetTelegramWebhook(cfg.Auth.Telegram)

	// Start the server
	r := router.New(cfg)
	router.RegisterRoutes(r, cfg, modelsFactory, annalistManager, bot)
	r.Logger.Fatal(router.Start(r, cfg))

}
