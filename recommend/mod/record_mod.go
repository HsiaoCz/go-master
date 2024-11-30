package mod

import (
	"context"

	"github.com/HsiaoCz/go-master/recommend/types"
	"gorm.io/gorm"
)

type RecordModInter interface {
	CreateRecord(context.Context, *types.Records) (*types.Records, error)
	GetRecordsByUserID(context.Context, string) ([]*types.Records, error)
}

type RecordMod struct {
	db *gorm.DB
}

func RecordModInit(db *gorm.DB) *RecordMod {
	return &RecordMod{
		db: db,
	}
}

func (r *RecordMod) CreateRecord(ctx context.Context, record *types.Records) (*types.Records, error) {
	tx := r.db.Debug().WithContext(ctx).Model(&types.Records{}).Create(record)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return record, nil
}

func (r *RecordMod) GetRecordsByUserID(ctx context.Context, user_id string) ([]*types.Records, error) {
	var records []*types.Records
	tx := r.db.Debug().WithContext(ctx).Model(&types.Records{}).Where("user_id = ?", user_id).Find(&records)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return records, nil
}
