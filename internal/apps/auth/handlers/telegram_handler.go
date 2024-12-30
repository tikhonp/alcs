package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/alcs/internal/util/auth"
)

func (ah *AuthHandler) TelegramCallback(c echo.Context) error {
	// Parse Telegram data from query params
	params := c.QueryParams()
	authData := make(map[string]string)
	for key, values := range params {
		authData[key] = values[0]
	}

	// Validate Telegram data
	if err := validateTelegramAuth(authData, ah.AuthCfg.Telegram.BotToken); err != nil {
        log.Println(err)
		return c.JSON(400, map[string]string{"error": "Invalid Telegram authentication"})
	}

	// Extract user data
	userId := authData["id"]
	username := authData["username"]
	photoUrl := authData["photo_url"]
	firstName := authData["first_name"]
	lastName := authData["last_name"]

	// Find or create a user in the database
	user, err := ah.Db.AuthUsers().FindOrCreateTelegramUser(userId, username, firstName, lastName, photoUrl)
	if err != nil {
		return err
	}

	// Log the user in and set the session
	auth.LoginByUserId(c, user.ID)

	// Redirect to the next path or the home page
	nextPath := c.QueryParam("next")
	if nextPath == "" {
		nextPath = "/"
	}
	return c.Redirect(302, nextPath)
}

// Validate Telegram authentication data
func validateTelegramAuth(authData map[string]string, botToken string) error {
	// Extract the Telegram `hash`
	hash, ok := authData["hash"]
	if !ok {
		return errors.New("hash is missing")
	}
	delete(authData, "hash")

	// Sort and concatenate the authData
	var dataCheckStrings []string
	for key, value := range authData {
		dataCheckStrings = append(dataCheckStrings, key+"="+value)
	}
	sort.Strings(dataCheckStrings)
	dataCheckString := strings.Join(dataCheckStrings, "\n")

	// Calculate the HMAC-SHA256 hash using the bot token
	secret := sha256.Sum256([]byte(botToken))
	h := hmac.New(sha256.New, secret[:])
	h.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	// Compare the provided hash with the calculated hash
	if hash != calculatedHash {
		return errors.New("data validation failed")
	}

	// Check if the authentication data is recent (10 minutes)
	timestamp, err := strconv.ParseInt(authData["auth_date"], 10, 64)
	if err != nil || time.Now().Unix()-timestamp > 600 {
		return errors.New("authentication data is outdated")
	}

	return nil
}
