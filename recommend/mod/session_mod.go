package mod

import (
	"context"

	"github.com/HsiaoCz/go-master/recommend/types"
	"gorm.io/gorm"
)

type SessionModInter interface {
	CreateSession(context.Context, *types.Sessions) (*types.Sessions, error)
	GetSessionByToken(context.Context, string) (*types.Sessions, error)
}

type SessionMod struct {
	db *gorm.DB
}

func SessionModInit(db *gorm.DB) *SessionMod {
	return &SessionMod{
		db: db,
	}
}

func (s *SessionMod) CreateSession(ctx context.Context, session *types.Sessions) (*types.Sessions, error) {
	tx := s.db.Debug().WithContext(ctx).Model(&types.Sessions{}).Create(session)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return session, nil
}

func (s *SessionMod)GetSessionByToken(ctx context.Context,token string)(*types.Sessions,error){
	var session types.Sessions
	tx:=s.db.Debug().WithContext(ctx).Model(&types.Sessions{}).Where("token = ?",token).First(&session)
	if tx.Error!=nil{
		return nil,tx.Error
	}
	return &session,nil
}