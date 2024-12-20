package strategy

import (
    "context"
    "crypto-quant/internal/models"
)

// Signal represents a trading signal
type Signal struct {
    Type      string  // "buy" or "sell"
    Price     float64
    Timestamp string
    Symbol    string
    Quantity  float64
}

// Strategy defines the interface for trading strategies
type Strategy interface {
    Initialize(ctx context.Context) error
    Update(ctx context.Context, kline models.Kline) error
    GetSignal(ctx context.Context) (*Signal, error)
}

// BaseStrategy provides common functionality for strategies
type BaseStrategy struct {
    Symbol    string
    Interval  string
    Klines    []models.Kline
    MaxKlines int
}

// Update adds a new kline and maintains the maximum number of klines
func (b *BaseStrategy) Update(_ context.Context, kline models.Kline) error {
    b.Klines = append(b.Klines, kline)
    if len(b.Klines) > b.MaxKlines {
        b.Klines = b.Klines[1:]
    }
    return nil
}
