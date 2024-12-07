package storage

import (
	"context"

	"github.com/HsiaoCz/go-master/stock/types"
	"gorm.io/gorm"
)

type SessionStoreInter interface {
	CreateSession(context.Context, *types.Sessions) (*types.Sessions, error)
	GetSessionByID(context.Context, string) (*types.Sessions, error)
	DeleteSessionByToken(context.Context, string) error
}

type SessionStore struct {
	db *gorm.DB
}

func SessionStoreInit(db *gorm.DB) *SessionStore {
	return &SessionStore{
		db: db,
	}
}

func (s *SessionStore) CreateSession(ctx context.Context, session *types.Sessions) (*types.Sessions, error) {
	return nil, nil
}

func (s *SessionStore) GetSessionByID(ctx context.Context, session_id string) (*types.Sessions, error) {
	return nil, nil
}

func (s *SessionStore) DeleteSessionByToken(ctx context.Context, token string) error {
	var session types.Sessions
	tx := s.db.Debug().WithContext(ctx).Model(&types.Sessions{}).Where("token = ?", token).Delete(&session)
	return tx.Error
}
