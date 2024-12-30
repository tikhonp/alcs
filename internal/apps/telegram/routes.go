package telegram

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/internal/apps/telegram/bot"
	"github.com/tikhonp/alcs/internal/apps/telegram/handlers"
	"github.com/tikhonp/alcs/internal/config"
	"github.com/tikhonp/alcs/internal/db"
)

func ConfigureTelegramPageGroup(g *echo.Group, cfg *config.Config, modelsFactory db.ModelsFactory, bot *bot.Bot) {
	th := handlers.TelegramHandler{Db: modelsFactory, Bot: bot}

	g.POST("/webhook", th.HandleTelegramWebhook)
}
