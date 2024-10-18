package main

import (
	"context"
	"fmt"

	"github.com/tikhonp/alcs/config"
	"github.com/tikhonp/alcs/db"
)

func main() {
	cfg, err := config.LoadFromPath(context.Background(), "config.pkl")
	if err != nil {
		panic(err)
	}
	fmt.Print(
		db.DataSourceName(cfg.Db),
	)
}

