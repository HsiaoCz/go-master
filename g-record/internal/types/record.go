package types

import "gorm.io/gorm"

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
