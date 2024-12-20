package strategy

import (
    "context"
    "crypto-quant/internal/indicators"
    "crypto-quant/internal/models"
)

// MACrossoverStrategy implements a simple moving average crossover strategy
type MACrossoverStrategy struct {
    BaseStrategy
    ShortPeriod int
    LongPeriod  int
    LastSignal  string
}

// NewMACrossoverStrategy creates a new MA crossover strategy
func NewMACrossoverStrategy(symbol string, shortPeriod, longPeriod int) *MACrossoverStrategy {
    return &MACrossoverStrategy{
        BaseStrategy: BaseStrategy{
            Symbol:    symbol,
            Interval:  "1h",
            MaxKlines: longPeriod + 10,
        },
        ShortPeriod: shortPeriod,
        LongPeriod:  longPeriod,
        LastSignal:  "",
    }
}

func (s *MACrossoverStrategy) Initialize(_ context.Context) error {
    s.Klines = make([]models.Kline, 0)
    return nil
}

func (s *MACrossoverStrategy) GetSignal(_ context.Context) (*Signal, error) {
    if len(s.Klines) < s.LongPeriod {
        return nil, nil
    }

    // Extract closing prices
    prices := make([]float64, len(s.Klines))
    for i, k := range s.Klines {
        prices[i] = k.Close
    }

    // Calculate MAs
    shortMA := indicators.CalculateMA(prices, s.ShortPeriod)
    longMA := indicators.CalculateMA(prices, s.LongPeriod)

    if len(shortMA) < 2 || len(longMA) < 2 {
        return nil, nil
    }

    // Check for crossover
    currentShort := shortMA[len(shortMA)-1]
    previousShort := shortMA[len(shortMA)-2]
    currentLong := longMA[len(longMA)-1]
    previousLong := longMA[len(longMA)-2]

    lastKline := s.Klines[len(s.Klines)-1]

    // Golden Cross (short MA crosses above long MA)
    if previousShort <= previousLong && currentShort > currentLong && s.LastSignal != "buy" {
        s.LastSignal = "buy"
        return &Signal{
            Type:      "buy",
            Price:     lastKline.Close,
            Timestamp: lastKline.CloseTime.String(),
            Symbol:    s.Symbol,
            Quantity:  1.0, // This should be calculated based on position sizing rules
        }, nil
    }

    // Death Cross (short MA crosses below long MA)
    if previousShort >= previousLong && currentShort < currentLong && s.LastSignal != "sell" {
        s.LastSignal = "sell"
        return &Signal{
            Type:      "sell",
            Price:     lastKline.Close,
            Timestamp: lastKline.CloseTime.String(),
            Symbol:    s.Symbol,
            Quantity:  1.0, // This should be calculated based on position sizing rules
        }, nil
    }

    return nil, nil
}
