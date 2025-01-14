package types

import "time"

type UserRelationships struct {
	ID               int64     `json:"id"`
	UserID           string    `json:"user_id"`
	RelatedUserID    string    `json:"related_user_id"`
	RelationshipType string    `json:"relationship_type"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
