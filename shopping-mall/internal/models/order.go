package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	ProductID uuid.UUID `json:"product_id" bson:"product_id"`
	Quantity  int       `json:"quantity" bson:"quantity"`
	Price     float64   `json:"price" bson:"price"`
}

type Order struct {
	ID            uuid.UUID    `json:"id" bson:"_id"`
	UserID        uuid.UUID    `json:"user_id" bson:"user_id"`
	Items         []OrderItem  `json:"items" bson:"items"`
	TotalAmount   float64      `json:"total_amount" bson:"total_amount"`
	Status        OrderStatus  `json:"status" bson:"status"`
	ShippingAddr  string       `json:"shipping_addr" bson:"shipping_addr"`
	PaymentMethod string       `json:"payment_method" bson:"payment_method"`
	CreatedAt     time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at" bson:"updated_at"`
}

type CreateOrderRequest struct {
	Items         []OrderItem `json:"items" binding:"required,min=1"`
	ShippingAddr  string      `json:"shipping_addr" binding:"required"`
	PaymentMethod string      `json:"payment_method" binding:"required"`
}
