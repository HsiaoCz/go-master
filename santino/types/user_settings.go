package types

import (
	"time"

	"github.com/google/uuid"
)

type UserSettings struct {
	SettingID               int64          `json:"setting_id"`
	UserID                  uuid.UUID      `json:"user_id"`
	PrivacyStatus           string         `json:"privacy_status"`
	NotificationPreferences map[string]any `json:"notification_preferences"`
	LanguagePreferences     string         `json:"language_preferences"`
	ThemePreferences        string         `json:"theme_preferences"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
}
