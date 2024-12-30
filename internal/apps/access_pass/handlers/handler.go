package handlers

import (
	"github.com/tikhonp/alcs/internal/apps/telegram/bot"
	"github.com/tikhonp/alcs/internal/db"
)

type AccessPassHandler struct {
	Db  db.ModelsFactory
	Bot *bot.Bot
}
