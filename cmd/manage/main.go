package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/tikhonp/alcs/config"
	"github.com/tikhonp/alcs/db"
)

type command string

const (
    // PrintDbString prints the database configurtion string
    // for sql connection
	PrintDbString command = "print-db-string"

    // CreateSuperUser creates a super user
    CreateSuperUser command = "create-super-user"
)

func (c *command) Set(value string) error {
	switch command(value) {
	case PrintDbString, CreateSuperUser:
		*c = command(value)
		return nil
	default:
		return fmt.Errorf("invalid command %s", value)
	}
}

func (c *command) String() string {
	return string(*c)
}

type manageConfig struct {
	command    command
	configPath string
}

func ParseFlags() *manageConfig {
	cfg := &manageConfig{}

    const commandUsage = "command to run. Available commands: print-db-string, create-super-user"
	flag.Var(&cfg.command, "command", commandUsage)
    flag.Var(&cfg.command, "c", commandUsage+" (shorthand)")

	flag.StringVar(&cfg.configPath, "config", "config.pkl", "path to config file")

	flag.Parse()

	return cfg
}

func printDbString(cfg *config.Config) {
	fmt.Print(
		db.DataSourceName(cfg.Db),
	)
}

func createSuperUser() {
    fmt.Println("Not implemented")
    os.Exit(1)
}

func main() {
	manageConfig := ParseFlags()

	cfg, err := config.LoadFromPath(context.Background(), manageConfig.configPath)
	if err != nil {
		panic(err)
	}

	switch manageConfig.command {
	case PrintDbString:
		printDbString(cfg)
    case CreateSuperUser:
        createSuperUser()        
	}
}
