package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/mail"
	"os"
	"syscall"

	"golang.org/x/term"

	"github.com/tikhonp/alcs/internal/config"
	"github.com/tikhonp/alcs/internal/db"
	"github.com/tikhonp/alcs/internal/db/models/auth"
	"github.com/tikhonp/alcs/internal/util/assert"
)

type command string

const (
	// PrintDbString prints the database configurtion string
	// for sql connection
	PrintDbString command = "print-db-string"

	// CreateSuperAdmin creates a super user
	CreateSuperAdmin command = "create-super-admin"
)

func (c *command) Set(value string) error {
	switch command(value) {
	case PrintDbString, CreateSuperAdmin:
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

func createSuperAdmin(users auth.Users) {
	var (
		email string
        password string
        firstName string
        lastName string
	)

	fmt.Print("Email adress: ")
	_, err := fmt.Scanln(&email)
	if err != nil {
		assert.NoError(err, "Scaning string should work???")
	}
	mail, err := mail.ParseAddress(email)
	if err != nil {
		fmt.Printf("Invalid email format: %v\n", err.Error())
		os.Exit(1)
	}
    email = mail.Address
    
    fmt.Print("Password: ")
    bytepw, err := term.ReadPassword(int(syscall.Stdin))
    fmt.Print("Repeat password: ")
    bytepwr, err := term.ReadPassword(int(syscall.Stdin))
    if !bytes.Equal(bytepw, bytepwr) {
        fmt.Println("Error: passwords does not match.")
        os.Exit(1)
    }
    password = string(bytepw)

    fmt.Print("First name (leave blank for null): ")
    _, err = fmt.Scanln(&firstName)
	if err != nil {
		assert.NoError(err, "Scaning string should work???")
	}

    fmt.Print("Last name (leave blank for null): ")
    _, err = fmt.Scanln(&lastName)
	if err != nil {
		assert.NoError(err, "Scaning string should work???")
	}

    err = users.CreateSuperAdmin(email, password, firstName, lastName)
    if err != nil {
        assert.NoError(err, "Superuser user creation failed")
    }
}

func connectToDb(cfg *config.Config) db.ModelsFactory {
	modelsFactory, err := db.Connect(cfg.Db)
	if err != nil {
        assert.NoError(err, "Failed database connection")
	}
    return modelsFactory
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
	case CreateSuperAdmin:
		createSuperAdmin(connectToDb(cfg).AuthUsers())
	}
}
