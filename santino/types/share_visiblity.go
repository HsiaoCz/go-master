package types

import "github.com/google/uuid"

type ShareVisibility struct {
	VisibilityID    int64     `json:"visibility_id"`
	ShareID         uuid.UUID `json:"share_id"`
	VisibleToUserID uuid.UUID `json:"visible_to_user_id"`
}
