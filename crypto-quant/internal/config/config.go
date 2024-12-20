package config

import (
    "encoding/json"
    "os"
    "time"
)

type Config struct {
    Exchange ExchangeConfig `json:"exchange"`
    Trading  TradingConfig  `json:"trading"`
    Risk     RiskConfig     `json:"risk"`
    Logger   LoggerConfig   `json:"logger"`
}

type ExchangeConfig struct {
    Name      string `json:"name"`
    APIKey    string `json:"apiKey"`
    APISecret string `json:"apiSecret"`
    TestNet   bool   `json:"testnet"`
}

type TradingConfig struct {
    Symbols        []string      `json:"symbols"`
    Interval       string        `json:"interval"`
    Strategy       string        `json:"strategy"`
    UpdateInterval time.Duration `json:"updateInterval"`
    Strategies     map[string]json.RawMessage `json:"strategies"`
}

type RiskConfig struct {
    InitialBalance float64 `json:"initialBalance"`
    RiskPerTrade   float64 `json:"riskPerTrade"`
    MaxDrawdown    float64 `json:"maxDrawdown"`
    MaxPositions   int     `json:"maxPositions"`
    DailyLossLimit float64 `json:"dailyLossLimit"`
}

type LoggerConfig struct {
    Level      string `json:"level"`
    File       string `json:"file"`
    MaxSize    int    `json:"maxSize"`    // megabytes
    MaxBackups int    `json:"maxBackups"` // number of backups
    MaxAge     int    `json:"maxAge"`     // days
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(path string) (*Config, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var config Config
    if err := json.NewDecoder(file).Decode(&config); err != nil {
        return nil, err
    }

    return &config, nil
}

// SaveConfig saves configuration to a JSON file
func SaveConfig(config *Config, path string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "    ")
    return encoder.Encode(config)
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
    return &Config{
        Exchange: ExchangeConfig{
            Name:    "binance",
            TestNet: true,
        },
        Trading: TradingConfig{
            Symbols:        []string{"BTC-USDT"},
            Interval:       "1h",
            Strategy:      "ma_crossover",
            UpdateInterval: time.Minute,
            Strategies: map[string]json.RawMessage{
                "ma_crossover": json.RawMessage(`{
                    "shortPeriod": 10,
                    "longPeriod": 20
                }`),
            },
        },
        Risk: RiskConfig{
            InitialBalance: 10000,
            RiskPerTrade:   1.0,
            MaxDrawdown:    20.0,
            MaxPositions:   3,
            DailyLossLimit: 5.0,
        },
        Logger: LoggerConfig{
            Level:      "info",
            File:       "trading.log",
            MaxSize:    100,
            MaxBackups: 3,
            MaxAge:     28,
        },
    }
}
