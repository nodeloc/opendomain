package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"opendomain/internal/config"
	"opendomain/pkg/timeutil"
)

type TelegramService struct {
	cfg *config.Config
}

func NewTelegramService(cfg *config.Config) *TelegramService {
	return &TelegramService{
		cfg: cfg,
	}
}

type telegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// SendHealthAlert sends a health alert notification to Telegram channel
func (s *TelegramService) SendHealthAlert(domain string, issues []string, action string) error {
	// Skip if Telegram is not configured
	if s.cfg.Telegram.BotToken == "" || s.cfg.Telegram.ChannelID == "" {
		return nil
	}

	// Build message
	message := fmt.Sprintf("🚨 *Domain Health Alert*\n\n")
	message += fmt.Sprintf("📍 Domain: `%s`\n", domain)
	message += fmt.Sprintf("⏰ Time: %s\n\n", timeutil.Now().UTC().Format("2006-01-02 15:04:05 UTC"))

	if len(issues) > 0 {
		message += "*Issues Detected:*\n"
		for _, issue := range issues {
			message += fmt.Sprintf("• %s\n", issue)
		}
		message += "\n"
	}

	if action != "" {
		message += fmt.Sprintf("⚠️ *Action: %s*\n", action)
	}

	return s.sendMessage(message)
}

// SendAutoSuspendNotification sends notification when domain is auto-suspended
func (s *TelegramService) SendAutoSuspendNotification(domain, reason string) error {
	if s.cfg.Telegram.BotToken == "" || s.cfg.Telegram.ChannelID == "" {
		return nil
	}

	message := fmt.Sprintf("🔒 *Domain Auto-Suspended*\n\n")
	message += fmt.Sprintf("📍 Domain: `%s`\n", domain)
	message += fmt.Sprintf("📝 Reason: %s\n", reason)
	message += fmt.Sprintf("⏰ Time: %s\n", timeutil.Now().UTC().Format("2006-01-02 15:04:05 UTC"))

	return s.sendMessage(message)
}

// SendDeletionWarning sends warning before domain deletion
func (s *TelegramService) SendDeletionWarning(domain string, daysRemaining int) error {
	if s.cfg.Telegram.BotToken == "" || s.cfg.Telegram.ChannelID == "" {
		return nil
	}

	message := fmt.Sprintf("⚠️ *Domain Deletion Warning*\n\n")
	message += fmt.Sprintf("📍 Domain: `%s`\n", domain)
	message += fmt.Sprintf("⏳ Days until deletion: *%d*\n", daysRemaining)
	message += fmt.Sprintf("⏰ Time: %s\n", timeutil.Now().UTC().Format("2006-01-02 15:04:05 UTC"))

	return s.sendMessage(message)
}

func (s *TelegramService) sendMessage(text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.cfg.Telegram.BotToken)

	msg := telegramMessage{
		ChatID:    s.cfg.Telegram.ChannelID,
		Text:      text,
		ParseMode: "Markdown",
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal telegram message: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}

	return nil
}
