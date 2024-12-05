package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Records struct {
	gorm.Model
	RecordID   string `gorm:"column:record_id;" json:"record_id"`
	UserID     string `gorm:"column:user_id;" json:"user_id"`
	BookID     string `gorm:"column:book_id;" json:"book_id"`
	Title      string `gorm:"column:title;" json:"title"`
	CoverImage string `gorm:"column:cover_image;" json:"cover_image"`
	Auther     string `gorm:"column:auther;" json:"auther"`
	TypeName   string `gorm:"column:type_name;" json:"type_name"`
	// 用户设备信息
	Device string `gorm:"column:device;" json:"device"`
}

type CreateRecordParams struct {
	UserID     string `gorm:"column:user_id;" json:"user_id"`
	BookID     string `gorm:"column:book_id;" json:"book_id"`
	Title      string `gorm:"column:title;" json:"title"`
	CoverImage string `gorm:"column:cover_image;" json:"cover_image"`
	Auther     string `gorm:"column:auther;" json:"auther"`
	TypeName   string `gorm:"column:type_name;" json:"type_name"`
	// 用户设备信息
	Device string `gorm:"column:device;" json:"device"`
}

func CreateRecordFromParams(record_params CreateRecordParams) *Records {
	return &Records{
		RecordID:   uuid.New().String(),
		UserID:     record_params.UserID,
		BookID:     record_params.BookID,
		Title:      record_params.Title,
		CoverImage: record_params.CoverImage,
		Auther:     record_params.Auther,
		TypeName:   record_params.TypeName,
		Device:     record_params.Device,
	}
}
