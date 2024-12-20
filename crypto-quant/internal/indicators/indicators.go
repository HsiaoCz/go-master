package indicators

// CalculateMA calculates Moving Average
func CalculateMA(prices []float64, period int) []float64 {
    if len(prices) < period {
        return nil
    }

    ma := make([]float64, len(prices)-period+1)
    for i := 0; i <= len(prices)-period; i++ {
        sum := 0.0
        for j := 0; j < period; j++ {
            sum += prices[i+j]
        }
        ma[i] = sum / float64(period)
    }
    return ma
}

// CalculateRSI calculates Relative Strength Index
func CalculateRSI(prices []float64, period int) []float64 {
    if len(prices) < period+1 {
        return nil
    }

    rsi := make([]float64, len(prices)-period)
    gains := make([]float64, 0)
    losses := make([]float64, 0)

    // Calculate price changes
    for i := 1; i < len(prices); i++ {
        change := prices[i] - prices[i-1]
        if change > 0 {
            gains = append(gains, change)
            losses = append(losses, 0)
        } else {
            gains = append(gains, 0)
            losses = append(losses, -change)
        }
    }

    // Calculate RSI
    for i := 0; i <= len(gains)-period; i++ {
        avgGain := average(gains[i : i+period])
        avgLoss := average(losses[i : i+period])
        
        if avgLoss == 0 {
            rsi[i] = 100
        } else {
            rs := avgGain / avgLoss
            rsi[i] = 100 - (100 / (1 + rs))
        }
    }
    return rsi
}

// Helper function to calculate average
func average(nums []float64) float64 {
    if len(nums) == 0 {
        return 0
    }
    sum := 0.0
    for _, n := range nums {
        sum += n
    }
    return sum / float64(len(nums))
}
