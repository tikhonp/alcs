package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tikhonp/alcs/internal/config"
	"github.com/tikhonp/alcs/internal/db/models/alcs"
)

type Bot struct {
	API *tgbotapi.BotAPI
}

// NewBot initializes the Telegram bot
func NewBot(botToken string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}
	return &Bot{API: api}, nil
}

// NotifyApproverWithInlineButtons sends a notification with inline buttons for approval or rejection
func (b *Bot) NotifyApproverWithInlineButtons(pass *alcs.AccessPass, accessPasses alcs.AccessPasses) error {
	messageText := fmt.Sprintf(
		"New Access Pass Request:\n\nVisitor: %s\nPurpose: %s\nValid From: %s\nValid Until: %s\n\nApprove or Reject this request below.",
		pass.VisitorName.String, pass.Purpose.String, pass.ValidFrom.Format("2006-01-02 15:04"), pass.ValidUntil.Format("2006-01-02 15:04"),
	)

	// Inline keyboard with Approve/Reject buttons
	approveButton := tgbotapi.NewInlineKeyboardButtonData("Approve", fmt.Sprintf("approve:%d", pass.ID))
	rejectButton := tgbotapi.NewInlineKeyboardButtonData("Reject", fmt.Sprintf("reject:%d", pass.ID))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(approveButton, rejectButton))

	// Send the message
	approverTelegramID, err := accessPasses.GetApproverTelegramID(pass.ID)
	if err != nil {
		log.Printf("Failed to get approver Telegram ID: %v", err)
		return err
	}
	msg := tgbotapi.NewMessage(approverTelegramID, messageText)
	msg.ReplyMarkup = keyboard
	_, err = b.API.Send(msg)
	if err != nil {
		log.Printf("Failed to send approval message to Telegram user %d: %v", approverTelegramID, err)
		return err
	}

	log.Printf("Approval message sent to Telegram user %d", approverTelegramID)
	return nil
}

// HandleUpdate processes Telegram webhook updates
func (b *Bot) HandleUpdate(update tgbotapi.Update, accessPasses alcs.AccessPasses) {
	if update.CallbackQuery != nil {
		callbackQuery := update.CallbackQuery
		callbackData := callbackQuery.Data

		// Parse the action and pass ID
		action, passIDStr, found := strings.Cut(callbackData, ":")
		if !found {
			callback := tgbotapi.NewCallback(callbackQuery.ID, "Invalid action")
			b.API.Send(callback)
			return
		}

		passID, err := strconv.Atoi(passIDStr)
		if err != nil {
			callback := tgbotapi.NewCallback(callbackQuery.ID, "Invalid pass ID")
			b.API.Send(callback)
			return
		}

		var responseMessage string
		switch action {
		case "approve":
			// Update the status to Approved
			err = accessPasses.UpdateAccessPassStatus(passID, alcs.AccessPassStatusApproved)
			if err == nil {
				responseMessage = fmt.Sprintf("Access pass #%d has been approved.", passID)
			} else {
				responseMessage = "Failed to approve the access pass."
			}

		case "reject":
			// Update the status to Rejected
			err = accessPasses.UpdateAccessPassStatus(passID, alcs.AccessPassStatusRejected)
			if err == nil {
				responseMessage = fmt.Sprintf("Access pass #%d has been rejected.", passID)
			} else {
				responseMessage = "Failed to reject the access pass."
			}

		default:
			responseMessage = "Unknown action."
		}

		// Respond to the callback query
		callback := tgbotapi.NewCallback(callbackQuery.ID, responseMessage)
		if _, err := b.API.Send(callback); err != nil {
			log.Printf("Failed to send callback response: %v", err)
		}

		// Optionally, edit the message to show the decision
		editMessage := tgbotapi.NewEditMessageText(
			callbackQuery.Message.Chat.ID,
			callbackQuery.Message.MessageID,
			fmt.Sprintf("Decision: %s\n\n%s", action, responseMessage),
		)
		b.API.Send(editMessage)
	}
}

// Set Telegram webhook
func (b *Bot) SetTelegramWebhook(cfg *config.TelegramAuth) error {
    webhookConfig, err := tgbotapi.NewWebhook(cfg.TelegramWebhookUrl)
	if err != nil {
		log.Fatalf("failed to create telegram webhook: %v", err)
		return err
	}
    _, err = b.API.Request(webhookConfig)
	if err != nil {
		log.Fatalf("failed to set telegram webhook: %v", err)
		return err
	}
	return nil
}
