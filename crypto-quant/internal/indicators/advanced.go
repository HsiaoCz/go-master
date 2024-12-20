package indicators

import "math"

// CalculateMACD calculates Moving Average Convergence Divergence
func CalculateMACD(prices []float64, fastPeriod, slowPeriod, signalPeriod int) ([]float64, []float64, []float64) {
    fastEMA := CalculateEMA(prices, fastPeriod)
    slowEMA := CalculateEMA(prices, slowPeriod)

    // Calculate MACD line
    macdLine := make([]float64, len(fastEMA))
    for i := 0; i < len(fastEMA); i++ {
        macdLine[i] = fastEMA[i] - slowEMA[i]
    }

    // Calculate Signal line (EMA of MACD line)
    signalLine := CalculateEMA(macdLine, signalPeriod)

    // Calculate histogram
    histogram := make([]float64, len(signalLine))
    for i := 0; i < len(signalLine); i++ {
        histogram[i] = macdLine[i] - signalLine[i]
    }

    return macdLine, signalLine, histogram
}

// CalculateEMA calculates Exponential Moving Average
func CalculateEMA(prices []float64, period int) []float64 {
    if len(prices) < period {
        return nil
    }

    multiplier := 2.0 / float64(period+1)
    ema := make([]float64, len(prices))
    
    // First EMA is calculated as simple moving average
    sum := 0.0
    for i := 0; i < period; i++ {
        sum += prices[i]
    }
    ema[period-1] = sum / float64(period)

    // Calculate EMA for remaining prices
    for i := period; i < len(prices); i++ {
        ema[i] = (prices[i]-ema[i-1])*multiplier + ema[i-1]
    }

    return ema
}

// BollingerBands calculates Bollinger Bands
func BollingerBands(prices []float64, period int, deviations float64) ([]float64, []float64, []float64) {
    if len(prices) < period {
        return nil, nil, nil
    }

    sma := CalculateMA(prices, period)
    upper := make([]float64, len(sma))
    lower := make([]float64, len(sma))

    for i := 0; i < len(sma); i++ {
        // Calculate standard deviation
        sum := 0.0
        for j := i; j < i+period && j < len(prices); j++ {
            diff := prices[j] - sma[i]
            sum += diff * diff
        }
        stdDev := math.Sqrt(sum / float64(period))

        upper[i] = sma[i] + deviations*stdDev
        lower[i] = sma[i] - deviations*stdDev
    }

    return upper, sma, lower
}

// CalculateATR calculates Average True Range
func CalculateATR(highs, lows, closes []float64, period int) []float64 {
    if len(highs) < 2 || len(lows) < 2 || len(closes) < 2 {
        return nil
    }

    trueRanges := make([]float64, len(highs))
    
    // Calculate True Range for each period
    for i := 1; i < len(highs); i++ {
        high := highs[i]
        low := lows[i]
        prevClose := closes[i-1]

        tr1 := high - low
        tr2 := math.Abs(high - prevClose)
        tr3 := math.Abs(low - prevClose)

        trueRanges[i] = math.Max(math.Max(tr1, tr2), tr3)
    }

    // Calculate ATR using smoothed moving average
    atr := make([]float64, len(trueRanges))
    sum := 0.0
    for i := 1; i <= period; i++ {
        sum += trueRanges[i]
    }
    atr[period] = sum / float64(period)

    for i := period + 1; i < len(trueRanges); i++ {
        atr[i] = (atr[i-1]*(float64(period)-1) + trueRanges[i]) / float64(period)
    }

    return atr
}
