package exchange

import (
    "context"
    "crypto-quant/internal/models"
)

// Exchange defines the interface for cryptocurrency exchanges
type Exchange interface {
    // Market Data
    GetKlines(ctx context.Context, symbol string, interval string, limit int) ([]models.Kline, error)
    GetOrderBook(ctx context.Context, symbol string, limit int) (*models.OrderBook, error)
    GetTrades(ctx context.Context, symbol string, limit int) ([]models.Trade, error)

    // Trading
    PlaceOrder(ctx context.Context, symbol string, side string, orderType string, quantity float64, price float64) (string, error)
    CancelOrder(ctx context.Context, symbol string, orderID string) error
    GetOrderStatus(ctx context.Context, symbol string, orderID string) (string, error)

    // Account
    GetBalance(ctx context.Context) (map[string]float64, error)
}

// Config holds exchange configuration
type Config struct {
    APIKey    string
    APISecret string
    BaseURL   string
}
