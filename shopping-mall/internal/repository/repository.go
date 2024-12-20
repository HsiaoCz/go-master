package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopping-mall/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *models.Product) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	ListProducts(ctx context.Context, skip, limit int) ([]*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	ListOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Order, error)
	UpdateOrder(ctx context.Context, order *models.Order) error
}
