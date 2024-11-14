package mod

import (
	"context"

	"github.com/HsiaoCz/go-master/recommend/types"
	"gorm.io/gorm"
)

type RecordModInter interface {
	CreateRecord(context.Context, *types.Records) (*types.Records, error)
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
	tx := r.db.Debug().Model(&types.Records{}).Create(record)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return record, nil
}
