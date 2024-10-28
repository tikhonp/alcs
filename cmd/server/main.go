package main

import (
	"context"
	"log"
	"os"

	"github.com/tikhonp/alcs/config"
	"github.com/tikhonp/alcs/db"
	"github.com/tikhonp/alcs/router"
)

func main() {

	// Read the configuration from the pkl file
	cfg, err := config.LoadFromPath(context.Background(), "config.pkl")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Connect to the database
	_, err = db.Connect(cfg.Db)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Start the server
	r := router.New(cfg)
	router.RegisterRoutes(r)
	r.Logger.Fatal(router.Start(r, cfg))

}

