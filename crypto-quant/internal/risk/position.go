package risk

import (
    "crypto-quant/internal/models"
    "math"
)

// PositionSizer handles position sizing calculations
type PositionSizer struct {
    AccountBalance float64
    RiskPerTrade   float64 // percentage of account balance to risk per trade
    MaxDrawdown    float64 // maximum allowed drawdown percentage
}

// NewPositionSizer creates a new position sizer
func NewPositionSizer(balance float64, riskPerTrade, maxDrawdown float64) *PositionSizer {
    return &PositionSizer{
        AccountBalance: balance,
        RiskPerTrade:   riskPerTrade,
        MaxDrawdown:    maxDrawdown,
    }
}

// CalculatePosition calculates the position size based on risk parameters
func (p *PositionSizer) CalculatePosition(price float64, stopLoss float64) float64 {
    if stopLoss >= price {
        return 0
    }

    // Calculate risk amount in quote currency
    riskAmount := p.AccountBalance * (p.RiskPerTrade / 100)
    
    // Calculate position size based on risk and stop loss
    riskPerUnit := math.Abs(price - stopLoss)
    if riskPerUnit == 0 {
        return 0
    }

    return riskAmount / riskPerUnit
}

// RiskManager handles risk management decisions
type RiskManager struct {
    MaxDrawdown        float64
    MaxPositions      int
    InitialBalance    float64
    CurrentBalance    float64
    OpenPositions     int
    DailyLossLimit    float64
    CurrentDailyLoss  float64
}

func NewRiskManager(initialBalance float64, maxDrawdown float64, maxPositions int, dailyLossLimit float64) *RiskManager {
    return &RiskManager{
        MaxDrawdown:      maxDrawdown,
        MaxPositions:    maxPositions,
        InitialBalance:  initialBalance,
        CurrentBalance:  initialBalance,
        DailyLossLimit:  dailyLossLimit,
    }
}

// CanTrade checks if a new trade can be taken based on risk parameters
func (r *RiskManager) CanTrade() bool {
    // Check drawdown
    drawdown := (r.InitialBalance - r.CurrentBalance) / r.InitialBalance * 100
    if drawdown >= r.MaxDrawdown {
        return false
    }

    // Check position limit
    if r.OpenPositions >= r.MaxPositions {
        return false
    }

    // Check daily loss limit
    if r.CurrentDailyLoss >= r.DailyLossLimit {
        return false
    }

    return true
}

// UpdateBalance updates the current balance and checks risk limits
func (r *RiskManager) UpdateBalance(newBalance float64) {
    r.CurrentBalance = newBalance
    r.CurrentDailyLoss = math.Max(0, r.InitialBalance-newBalance)
}

// TradeExit represents a trade exit signal with reason
type TradeExit struct {
    Reason string
    Price  float64
}

// ExitChecker checks for trade exit conditions
type ExitChecker struct {
    TrailingStop     float64 // percentage for trailing stop
    MaxLoss          float64 // maximum loss percentage per trade
    ProfitTarget     float64 // profit target percentage
    TimeStop         int     // maximum time in position (in minutes)
}

func (e *ExitChecker) CheckExit(position models.Position, currentPrice float64, timeInPosition int) *TradeExit {
    entryPrice := position.EntryPrice
    
    // Check stop loss
    if currentPrice <= position.StopLoss {
        return &TradeExit{
            Reason: "Stop Loss Hit",
            Price:  position.StopLoss,
        }
    }

    // Check profit target
    profitTarget := entryPrice * (1 + e.ProfitTarget/100)
    if currentPrice >= profitTarget {
        return &TradeExit{
            Reason: "Profit Target Hit",
            Price:  profitTarget,
        }
    }

    // Check trailing stop
    if position.TrailingStop > 0 && currentPrice < position.TrailingStop {
        return &TradeExit{
            Reason: "Trailing Stop Hit",
            Price:  position.TrailingStop,
        }
    }

    // Check time stop
    if timeInPosition >= e.TimeStop {
        return &TradeExit{
            Reason: "Time Stop Hit",
            Price:  currentPrice,
        }
    }

    return nil
}
