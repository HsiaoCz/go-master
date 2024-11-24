package mod

import (
	"context"

	"github.com/HsiaoCz/go-master/recommend/types"
	"gorm.io/gorm"
)

type BookModInter interface {
	CreateBook(context.Context, *types.Books) (*types.Books, error)
	GetBookByAuther(context.Context, string) ([]*types.Books, error)
}

type BookMod struct {
	db *gorm.DB
}

func BookModInit(db *gorm.DB) *BookMod {
	return &BookMod{
		db: db,
	}
}

func (b *BookMod) CreateBook(ctx context.Context, book *types.Books) (*types.Books, error) {
	tx := b.db.Debug().WithContext(ctx).Model(&types.Books{}).Create(book)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return book, nil
}

func (b *BookMod) GetBookByAuther(ctx context.Context, auther_name string) ([]*types.Books, error) {
	return nil, nil
}
