package handlers

import (
	"github.com/tikhonp/alcs/internal/db"
	"github.com/tikhonp/alcs/internal/util/annalist"
)

type AuthHandler struct {
	Db       db.ModelsFactory
	Annalist annalist.Annalist
}
