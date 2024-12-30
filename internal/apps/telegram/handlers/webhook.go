package handlers

import (
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/labstack/echo/v4"
)

func (th *TelegramHandler) HandleTelegramWebhook(c echo.Context) error {
	var update tgbotapi.Update
	if err := c.Bind(&update); err != nil {
		log.Printf("Failed to bind Telegram update: %v", err)
		return err
	}
	th.Bot.HandleUpdate(update, th.Db.AlcsAccessPasses())
	return c.NoContent(http.StatusOK)
}
