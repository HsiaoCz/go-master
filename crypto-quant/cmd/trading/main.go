package main

import (
	"context"
	"crypto-quant/internal/config"
	"crypto-quant/internal/exchange/binance"
	"crypto-quant/internal/risk"
	"crypto-quant/internal/strategy"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		// If config doesn't exist, create default
		cfg = config.DefaultConfig()
		if err := config.SaveConfig(cfg, "config.json"); err != nil {
			log.Fatalf("Failed to save default config: %v", err)
		}
	}

	// Initialize exchange client
	client := binance.NewClient(cfg.Exchange.APIKey, cfg.Exchange.APISecret)

	// Initialize risk management
	riskManager := risk.NewRiskManager(
		cfg.Risk.InitialBalance,
		cfg.Risk.MaxDrawdown,
		cfg.Risk.MaxPositions,
		cfg.Risk.DailyLossLimit,
	)

	positionSizer := risk.NewPositionSizer(
		cfg.Risk.InitialBalance,
		cfg.Risk.RiskPerTrade,
		cfg.Risk.MaxDrawdown,
	)

	// Initialize strategy
	strat := strategy.NewMACrossoverStrategy(cfg.Trading.Symbols[0], 10, 20)
	if err := strat.Initialize(context.Background()); err != nil {
		log.Fatalf("Failed to initialize strategy: %v", err)
	}

	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Create ticker for the trading loop
	ticker := time.NewTicker(cfg.Trading.UpdateInterval)
	defer ticker.Stop()

	log.Println("Starting trading bot...")

	// Main trading loop
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down gracefully...")
			return
		case <-ticker.C:
			// Fetch latest market data
			klines, err := client.GetKlines(ctx, cfg.Trading.Symbols[0], cfg.Trading.Interval, 100)
			if err != nil {
				log.Printf("Error fetching klines: %v", err)
				continue
			}

			// Update strategy with new data
			if err := strat.Update(ctx, klines[len(klines)-1]); err != nil {
				log.Printf("Error updating strategy: %v", err)
				continue
			}

			// Check for trading signals
			if !riskManager.CanTrade() {
				continue
			}

			signal, err := strat.GetSignal(ctx)
			if err != nil {
				log.Printf("Error getting signal: %v", err)
				continue
			}

			if signal != nil {
				// Calculate position size
				orderBook, err := client.GetOrderBook(ctx, cfg.Trading.Symbols[0], 5)
				if err != nil {
					log.Printf("Error fetching order book: %v", err)
					continue
				}

				var stopLoss float64
				if signal.Type == "buy" {
					stopLoss = orderBook.Bids[0].Price * 0.99 // 1% below current bid
				} else {
					stopLoss = orderBook.Asks[0].Price * 1.01 // 1% above current ask
				}

				size := positionSizer.CalculatePosition(signal.Price, stopLoss)
				if size <= 0 {
					continue
				}

				// Place order
				orderID, err := client.PlaceOrder(ctx, signal.Symbol, signal.Type, "MARKET", size, 0)
				if err != nil {
					log.Printf("Error placing order: %v", err)
					continue
				}

				log.Printf("Placed %s order: %s, Size: %f, Price: %f", signal.Type, orderID, size, signal.Price)
			}
		}
	}
}
