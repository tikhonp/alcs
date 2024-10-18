package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tikhonp/alcs/config"
	"github.com/tikhonp/alcs/db"
)

func main() {
	cfg, err := config.LoadFromPath(context.Background(), "config.pkl")
	if err != nil {
		panic(err)
	}
	_, err = db.Connect(cfg.Db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Lol")
	for {
		time.Sleep(1 * time.Hour)
	}
}

