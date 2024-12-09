package storage

import (
	"context"

	"github.com/HsiaoCz/go-master/recommend/types"
	"gorm.io/gorm"
)

type BookStorer interface {
	CreateBook(context.Context, *types.Books) (*types.Books, error)
	GetBookByAuther(context.Context, string) ([]*types.Books, error)
	GetBookByID(context.Context, string) (*types.Books, error)
}

type BookStore struct {
	db *gorm.DB
}

func BookStoreInit(db *gorm.DB) *BookStore {
	return &BookStore{
		db: db,
	}
}

func (b *BookStore) CreateBook(ctx context.Context, book *types.Books) (*types.Books, error) {
	tx := b.db.Debug().WithContext(ctx).Model(&types.Books{}).Create(book)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return book, nil
}

func (b *BookStore) GetBookByAuther(ctx context.Context, auther_name string) ([]*types.Books, error) {
	var books []*types.Books
	tx := b.db.Debug().WithContext(ctx).Model(&types.Books{}).Where("auther = ?", auther_name).Find(&books)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return books, nil
}

func (b *BookStore) GetBookByID(ctx context.Context, book_id string) (*types.Books, error) {
	var book types.Books
	tx := b.db.Debug().WithContext(ctx).Model(&types.Books{}).Where("book_id = ?", book_id).Find(&book)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &book, nil
}
