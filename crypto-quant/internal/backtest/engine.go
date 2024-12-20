package backtest

import (
    "crypto-quant/internal/models"
    "crypto-quant/internal/strategy"
    "crypto-quant/internal/risk"
    "time"
)

// Trade represents a backtesting trade
type Trade struct {
    Symbol    string
    Side      string
    EntryTime time.Time
    ExitTime  time.Time
    EntryPrice float64
    ExitPrice  float64
    Quantity   float64
    PnL        float64
    PnLPercent float64
}

// BacktestResult contains the results of a backtest
type BacktestResult struct {
    Trades            []Trade
    TotalPnL         float64
    WinRate          float64
    MaxDrawdown      float64
    SharpeRatio      float64
    StartBalance     float64
    EndBalance       float64
    TotalTrades      int
    WinningTrades    int
    LosingTrades     int
    AveragePnL       float64
    AverageWinPnL    float64
    AverageLossPnL   float64
    LargestWin       float64
    LargestLoss      float64
}

// Engine is the backtesting engine
type Engine struct {
    strategy  strategy.Strategy
    risk      *risk.RiskManager
    sizer     *risk.PositionSizer
    data      []models.Kline
    balance   float64
    positions []models.Position
    trades    []Trade
}

func NewEngine(strat strategy.Strategy, initialBalance float64) *Engine {
    return &Engine{
        strategy:  strat,
        risk:     risk.NewRiskManager(initialBalance, 20, 3, initialBalance*0.05),
        sizer:    risk.NewPositionSizer(initialBalance, 1.0, 20),
        balance:  initialBalance,
    }
}

// LoadData loads historical data for backtesting
func (e *Engine) LoadData(data []models.Kline) {
    e.data = data
}

// Run executes the backtest
func (e *Engine) Run() *BacktestResult {
    var trades []Trade
    maxBalance := e.balance
    minDrawdown := 0.0

    for i := range e.data {
        // Update strategy with new data
        e.strategy.Update(nil, e.data[i])

        // Check for signals
        signal, err := e.strategy.GetSignal(nil)
        if err != nil {
            continue
        }

        if signal != nil {
            // Handle existing positions
            for j, pos := range e.positions {
                if signal.Type == "sell" && pos.Side == "long" || signal.Type == "buy" && pos.Side == "short" {
                    trade := e.closePosition(pos, signal.Price, e.data[i].CloseTime)
                    trades = append(trades, trade)
                    // Remove the position
                    e.positions = append(e.positions[:j], e.positions[j+1:]...)
                }
            }

            // Handle new positions
            if len(e.positions) == 0 {
                // Create new position logic here
                // You may want to add position sizing and risk management here
                newPos := models.Position{
                    Symbol:     signal.Symbol,
                    Side:       signal.Type,
                    EntryPrice: signal.Price,
                    Quantity:   signal.Quantity,
                    EntryTime:  e.data[i].CloseTime,
                }
                e.positions = append(e.positions, newPos)
            }
        }

        // Track maximum drawdown
        if e.balance > maxBalance {
            maxBalance = e.balance
        }
        drawdown := (maxBalance - e.balance) / maxBalance * 100
        if drawdown > minDrawdown {
            minDrawdown = drawdown
        }
    }

    // Calculate final statistics
    return e.calculateResults(minDrawdown, trades)
}

func (e *Engine) checkPositions(kline models.Kline) {
    var remainingPositions []models.Position

    for _, pos := range e.positions {
        // Check stop loss
        if kline.Low <= pos.StopLoss {
            trade := e.closePosition(pos, pos.StopLoss, kline.CloseTime)
            e.trades = append(e.trades, trade)
            continue
        }

        // Update trailing stop if applicable
        if pos.TrailingStop > 0 {
            newStop := kline.High * (1 - pos.TrailingStop/100)
            if newStop > pos.StopLoss {
                pos.StopLoss = newStop
            }
        }

        remainingPositions = append(remainingPositions, pos)
    }

    e.positions = remainingPositions
}

func (e *Engine) closePosition(pos models.Position, exitPrice float64, exitTime time.Time) Trade {
    pnl := (exitPrice - pos.EntryPrice) * pos.Quantity
    if pos.Side == "short" {
        pnl = (pos.EntryPrice - exitPrice) * pos.Quantity
    }
    
    pnlPercent := (pnl / (pos.EntryPrice * pos.Quantity)) * 100
    
    trade := Trade{
        Symbol:     pos.Symbol,
        Side:       pos.Side,
        EntryTime:  pos.EntryTime,
        ExitTime:   exitTime,
        EntryPrice: pos.EntryPrice,
        ExitPrice:  exitPrice,
        Quantity:   pos.Quantity,
        PnL:        pnl,
        PnLPercent: pnlPercent,
    }
    
    e.balance += pnl
    return trade
}

func (e *Engine) calculateStopLoss(signal *strategy.Signal, kline models.Kline) float64 {
    // Simple implementation - you might want to make this more sophisticated
    atr := kline.High - kline.Low // Simplified ATR
    if signal.Type == "buy" {
        return kline.Close - 2*atr
    }
    return kline.Close + 2*atr
}

func (e *Engine) calculateResults(maxDrawdown float64, trades []Trade) *BacktestResult {
    result := &BacktestResult{
        Trades:        trades,
        StartBalance: e.balance,
        EndBalance:   e.balance,
        MaxDrawdown:  maxDrawdown,
        TotalTrades:  len(trades),
    }

    var totalPnL, totalWinPnL, totalLossPnL float64
    for _, trade := range trades {
        totalPnL += trade.PnL
        if trade.PnL > 0 {
            result.WinningTrades++
            totalWinPnL += trade.PnL
            if trade.PnL > result.LargestWin {
                result.LargestWin = trade.PnL
            }
        } else {
            result.LosingTrades++
            totalLossPnL += trade.PnL
            if trade.PnL < result.LargestLoss {
                result.LargestLoss = trade.PnL
            }
        }
    }

    result.TotalPnL = totalPnL
    if result.TotalTrades > 0 {
        result.WinRate = float64(result.WinningTrades) / float64(result.TotalTrades)
        result.AveragePnL = totalPnL / float64(result.TotalTrades)
        if result.WinningTrades > 0 {
            result.AverageWinPnL = totalWinPnL / float64(result.WinningTrades)
        }
        if result.LosingTrades > 0 {
            result.AverageLossPnL = totalLossPnL / float64(result.LosingTrades)
        }
    }

    return result
}
