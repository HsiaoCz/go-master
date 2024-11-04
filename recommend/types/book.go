package types

import "gorm.io/gorm"

type Books struct {
	gorm.Model
	BookID string  `gorm:"column:book_id;" json:"book_id"`
	Auther string  `gorm:"column:auther;" json:"auther"`
	Title  string  `gorm:"column:title;" json:"title"`
	Price  float64 `gorm:"column:price;" json:"price"`
}
