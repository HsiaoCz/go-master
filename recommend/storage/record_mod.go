package storage

import (
	"context"

	"github.com/HsiaoCz/go-master/recommend/types"
	"gorm.io/gorm"
)

type RecordStorer interface {
	CreateRecord(context.Context, *types.Records) (*types.Records, error)
	GetRecordsByUserID(context.Context, string) ([]*types.Records, error)
}

type RecordStore struct {
	db *gorm.DB
}

func RecordStoreInit(db *gorm.DB) *RecordStore {
	return &RecordStore{
		db: db,
	}
}

func (r *RecordStore) CreateRecord(ctx context.Context, record *types.Records) (*types.Records, error) {
	tx := r.db.Debug().WithContext(ctx).Model(&types.Records{}).Create(record)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return record, nil
}

func (r *RecordStore) GetRecordsByUserID(ctx context.Context, user_id string) ([]*types.Records, error) {
	var records []*types.Records
	tx := r.db.Debug().WithContext(ctx).Model(&types.Records{}).Where("user_id = ?", user_id).Find(&records)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return records, nil
}
