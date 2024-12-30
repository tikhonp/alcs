package handlers

import (
	"github.com/tikhonp/alcs/internal/config"
	"github.com/tikhonp/alcs/internal/db"
	"github.com/tikhonp/alcs/internal/util/annalist"
)

type AuthHandler struct {
	AuthCfg  *config.Auth
	Db       db.ModelsFactory
	Annalist annalist.Annalist
}

