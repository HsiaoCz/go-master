package models

import "time"

// Kline represents candlestick data
type Kline struct {
    OpenTime  time.Time
    Open      float64
    High      float64
    Low       float64
    Close     float64
    Volume    float64
    CloseTime time.Time
}

// Trade represents a single trade
type Trade struct {
    ID        string
    Price     float64
    Amount    float64
    Side      string    // buy or sell
    Timestamp time.Time
}

// OrderBook represents market depth
type OrderBook struct {
    Timestamp time.Time
    Bids      []OrderBookEntry
    Asks      []OrderBookEntry
}

// OrderBookEntry represents a single entry in the order book
type OrderBookEntry struct {
    Price  float64
    Amount float64
}

// Position represents an open trading position
type Position struct {
    Symbol      string
    Side        string    // "long" or "short"
    EntryPrice  float64
    Quantity    float64
    EntryTime   time.Time
    StopLoss    float64
    TrailingStop float64
    TakeProfit  float64
}
