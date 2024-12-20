package binance

import (
    "context"
    "crypto-quant/internal/models"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strconv"
    "time"
)

const (
    baseURL = "https://api.binance.com"
)

type Client struct {
    apiKey     string
    apiSecret  string
    httpClient *http.Client
}

func NewClient(apiKey, apiSecret string) *Client {
    return &Client{
        apiKey:     apiKey,
        apiSecret:  apiSecret,
        httpClient: &http.Client{Timeout: 10 * time.Second},
    }
}

func (c *Client) sign(params url.Values) string {
    mac := hmac.New(sha256.New, []byte(c.apiSecret))
    mac.Write([]byte(params.Encode()))
    return hex.EncodeToString(mac.Sum(nil))
}

func (c *Client) doRequest(ctx context.Context, method, endpoint string, params url.Values, signed bool) ([]byte, error) {
    if params == nil {
        params = url.Values{}
    }

    if signed {
        timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
        params.Set("timestamp", timestamp)
        signature := c.sign(params)
        params.Set("signature", signature)
    }

    reqURL := baseURL + endpoint
    if len(params) > 0 {
        reqURL += "?" + params.Encode()
    }

    req, err := http.NewRequestWithContext(ctx, method, reqURL, nil)
    if err != nil {
        return nil, fmt.Errorf("creating request: %w", err)
    }

    if c.apiKey != "" {
        req.Header.Set("X-MBX-APIKEY", c.apiKey)
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("executing request: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("reading response: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API error: %s", string(body))
    }

    return body, nil
}

// GetKlines implements the Exchange interface for getting kline/candlestick data
func (c *Client) GetKlines(ctx context.Context, symbol string, interval string, limit int) ([]models.Kline, error) {
    params := url.Values{}
    params.Set("symbol", symbol)
    params.Set("interval", interval)
    params.Set("limit", strconv.Itoa(limit))

    data, err := c.doRequest(ctx, "GET", "/api/v3/klines", params, false)
    if err != nil {
        return nil, err
    }

    var rawKlines [][]interface{}
    if err := json.Unmarshal(data, &rawKlines); err != nil {
        return nil, fmt.Errorf("parsing klines: %w", err)
    }

    klines := make([]models.Kline, len(rawKlines))
    for i, raw := range rawKlines {
        openTime := time.Unix(int64(raw[0].(float64))/1000, 0)
        closeTime := time.Unix(int64(raw[6].(float64))/1000, 0)
        
        open, _ := strconv.ParseFloat(raw[1].(string), 64)
        high, _ := strconv.ParseFloat(raw[2].(string), 64)
        low, _ := strconv.ParseFloat(raw[3].(string), 64)
        close, _ := strconv.ParseFloat(raw[4].(string), 64)
        volume, _ := strconv.ParseFloat(raw[5].(string), 64)

        klines[i] = models.Kline{
            OpenTime:  openTime,
            Open:      open,
            High:      high,
            Low:       low,
            Close:     close,
            Volume:    volume,
            CloseTime: closeTime,
        }
    }

    return klines, nil
}

// GetOrderBook implements the Exchange interface for getting market depth
func (c *Client) GetOrderBook(ctx context.Context, symbol string, limit int) (*models.OrderBook, error) {
    params := url.Values{}
    params.Set("symbol", symbol)
    params.Set("limit", strconv.Itoa(limit))

    data, err := c.doRequest(ctx, "GET", "/api/v3/depth", params, false)
    if err != nil {
        return nil, err
    }

    var response struct {
        LastUpdateID int64      `json:"lastUpdateId"`
        Bids        [][]string `json:"bids"`
        Asks        [][]string `json:"asks"`
    }

    if err := json.Unmarshal(data, &response); err != nil {
        return nil, fmt.Errorf("parsing order book: %w", err)
    }

    orderBook := &models.OrderBook{
        Timestamp: time.Now(),
        Bids:      make([]models.OrderBookEntry, len(response.Bids)),
        Asks:      make([]models.OrderBookEntry, len(response.Asks)),
    }

    for i, bid := range response.Bids {
        price, _ := strconv.ParseFloat(bid[0], 64)
        amount, _ := strconv.ParseFloat(bid[1], 64)
        orderBook.Bids[i] = models.OrderBookEntry{Price: price, Amount: amount}
    }

    for i, ask := range response.Asks {
        price, _ := strconv.ParseFloat(ask[0], 64)
        amount, _ := strconv.ParseFloat(ask[1], 64)
        orderBook.Asks[i] = models.OrderBookEntry{Price: price, Amount: amount}
    }

    return orderBook, nil
}

// PlaceOrder implements the Exchange interface for placing orders
func (c *Client) PlaceOrder(ctx context.Context, symbol string, side string, orderType string, quantity float64, price float64) (string, error) {
    params := url.Values{}
    params.Set("symbol", symbol)
    params.Set("side", side)
    params.Set("type", orderType)
    params.Set("quantity", strconv.FormatFloat(quantity, 'f', -1, 64))
    
    if orderType == "LIMIT" {
        params.Set("price", strconv.FormatFloat(price, 'f', -1, 64))
        params.Set("timeInForce", "GTC")
    }

    data, err := c.doRequest(ctx, "POST", "/api/v3/order", params, true)
    if err != nil {
        return "", err
    }

    var response struct {
        OrderID int64 `json:"orderId"`
    }

    if err := json.Unmarshal(data, &response); err != nil {
        return "", fmt.Errorf("parsing order response: %w", err)
    }

    return strconv.FormatInt(response.OrderID, 10), nil
}

// Additional methods to implement the Exchange interface...
