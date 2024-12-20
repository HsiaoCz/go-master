package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Price       float64   `json:"price" bson:"price"`
	Category    string    `json:"category" bson:"category"`
	Stock       int       `json:"stock" bson:"stock"`
	ImageURL    string    `json:"image_url" bson:"image_url"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Category    string  `json:"category" binding:"required"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	ImageURL    string  `json:"image_url" binding:"required,url"`
}
