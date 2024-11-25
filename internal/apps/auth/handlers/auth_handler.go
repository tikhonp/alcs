package handlers

import "github.com/tikhonp/alcs/internal/db"

type AuthHandler struct {
	Db db.ModelsFactory
}
