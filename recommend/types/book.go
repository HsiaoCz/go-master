package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Books struct {
	gorm.Model
	BookID       string  `gorm:"column:book_id;" json:"book_id"`
	Auther       string  `gorm:"column:auther;" json:"auther"`
	Title        string  `gorm:"column:title;" json:"title"`
	Price        float64 `gorm:"column:price;" json:"price"`
	Stock        int     `gorm:"column:stock;" json:"stock"`
	CategoryName string  `gorm:"column:category_name;" json:"category_name"`
	Descriptions string  `gorm:"column:descriptions;" json:"descriptions"`
	CoverImage   string  `gorm:"column:cover_image;" json:"cover_image"`
}

type CreateBookParams struct {
	Auther       string  `json:"auther"`
	Title        string  `json:"title"`
	Price        float64 `json:"float64"`
	Stock        int     `json:"stock"`
	CatehoryName string  `json:"category_name"`
	Descriptions string  `json:"descriptions"`
	CoverImage   string  `json:"cover_image"`
}

func CreateBookFromParams(params CreateBookParams) *Books {
	return &Books{
		Auther:       params.Auther,
		Title:        params.Title,
		Price:        params.Price,
		Stock:        params.Stock,
		CategoryName: params.CatehoryName,
		Descriptions: params.Descriptions,
		CoverImage:   params.CoverImage,
		BookID:       uuid.New().String(),
	}
}
