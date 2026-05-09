package marketdata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var eastmoneyQuoteURL = "https://push2.eastmoney.com/api/qt/stock/get"
var tencentQuoteURL = "https://qt.gtimg.cn/q="

type eastmoneyResponse struct {
	Code int `json:"rc"`
	Data struct {
		Price         float64 `json:"f43"`
		High          float64 `json:"f44"`
		Low           float64 `json:"f45"`
		Open          float64 `json:"f46"`
		Volume        int64   `json:"f47"`
		Symbol        string  `json:"f57"`
		Name          string  `json:"f58"`
		PreviousClose float64 `json:"f60"`
		Timestamp     int64   `json:"f86"`
	} `json:"data"`
}

func FetchEastmoneyQuote(ctx context.Context, request QuoteRequest, client *http.Client) (Snapshot, error) {
	symbol := strings.ToUpper(strings.TrimSpace(request.Symbol))
	secid, err := eastmoneySecID(symbol)
	if err != nil {
		return Snapshot{}, err
	}
	u, err := url.Parse(eastmoneyQuoteURL)
	if err != nil {
		return Snapshot{}, err
	}
	query := u.Query()
	query.Set("secid", secid)
	query.Set("fields", "f43,f44,f45,f46,f47,f57,f58,f60,f86")
	u.RawQuery = query.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return Snapshot{}, err
	}
	httpReq.Header.Set("Referer", "https://finance.eastmoney.com/")
	httpReq.Header.Set("User-Agent", "jijin-stock-monitor/0.1")

	if client == nil {
		client = &http.Client{Timeout: 8 * time.Second}
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return Snapshot{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Snapshot{}, fmt.Errorf("eastmoney returned status %d", resp.StatusCode)
	}

	var parsed eastmoneyResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return Snapshot{}, err
	}
	if parsed.Code != 0 {
		return Snapshot{}, fmt.Errorf("eastmoney returned code %d", parsed.Code)
	}
	if parsed.Data.Price <= 0 {
		return Snapshot{}, fmt.Errorf("eastmoney quote missing price for %s", request.Key())
	}
	dataTime := time.Now().UTC()
	if parsed.Data.Timestamp > 0 {
		dataTime = time.Unix(parsed.Data.Timestamp, 0).UTC()
	}
	previousClose := scaleEastmoneyPrice(parsed.Data.PreviousClose)
	price := scaleEastmoneyPrice(parsed.Data.Price)
	changePercent := 0.0
	if previousClose > 0 {
		changePercent = (price - previousClose) / previousClose * 100
	}
	name := strings.TrimSpace(parsed.Data.Name)
	if name == "" {
		name = symbol
	}
	return Snapshot{
		Market:        "CN",
		Symbol:        symbol,
		Name:          name,
		Open:          scaleEastmoneyPrice(parsed.Data.Open),
		High:          scaleEastmoneyPrice(parsed.Data.High),
		Low:           scaleEastmoneyPrice(parsed.Data.Low),
		Price:         price,
		PreviousClose: previousClose,
		ChangePercent: changePercent,
		Volume:        parsed.Data.Volume * 100,
		Source:        "eastmoney",
		DataTime:      dataTime,
		CreatedAt:     time.Now().UTC(),
	}, nil
}

func FetchTencentQuote(ctx context.Context, request QuoteRequest, client *http.Client) (Snapshot, error) {
	symbol := strings.ToUpper(strings.TrimSpace(request.Symbol))
	code, err := tencentSymbol(symbol)
	if err != nil {
		return Snapshot{}, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, tencentQuoteURL+code, nil)
	if err != nil {
		return Snapshot{}, err
	}
	httpReq.Header.Set("Referer", "https://gu.qq.com/")
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 jijin-stock-monitor/0.1")
	if client == nil {
		client = &http.Client{Timeout: 8 * time.Second}
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return Snapshot{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Snapshot{}, fmt.Errorf("tencent returned status %d", resp.StatusCode)
	}

	// 腾讯行情返回 GBK 编码的波浪线分隔数据，用作东方财富异常时的备用行情源。
	body, err := io.ReadAll(transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return Snapshot{}, err
	}
	fields := strings.Split(strings.TrimSpace(string(body)), "~")
	if len(fields) < 35 {
		return Snapshot{}, fmt.Errorf("tencent quote missing fields for %s", request.Key())
	}
	price, err := parseFloat(fields[3])
	if err != nil || price <= 0 {
		return Snapshot{}, fmt.Errorf("tencent quote missing price for %s", request.Key())
	}
	previousClose, _ := parseFloat(fields[4])
	open, _ := parseFloat(fields[5])
	volumeHands, _ := parseFloat(fields[6])
	changePercent, _ := parseFloat(fields[32])
	high, _ := parseFloat(fields[33])
	low, _ := parseFloat(fields[34])
	dataTime := parseTencentTime(fields[30])
	name := strings.TrimSpace(fields[1])
	if name == "" {
		name = symbol
	}
	return Snapshot{
		Market:        "CN",
		Symbol:        symbol,
		Name:          name,
		Open:          open,
		High:          high,
		Low:           low,
		Price:         price,
		PreviousClose: previousClose,
		ChangePercent: changePercent,
		Volume:        int64(volumeHands * 100),
		Source:        "tencent",
		DataTime:      dataTime,
		CreatedAt:     time.Now().UTC(),
	}, nil
}

func eastmoneySecID(symbol string) (string, error) {
	if len(symbol) != 6 {
		return "", errors.New("CN stock symbol must be 6 digits")
	}
	switch symbol[0] {
	case '0', '2', '3':
		return "0." + symbol, nil
	case '6', '9':
		return "1." + symbol, nil
	default:
		return "", fmt.Errorf("unsupported CN stock prefix %q", symbol[:1])
	}
}

func tencentSymbol(symbol string) (string, error) {
	if len(symbol) != 6 {
		return "", errors.New("CN stock symbol must be 6 digits")
	}
	switch symbol[0] {
	case '0', '2', '3':
		return "sz" + symbol, nil
	case '6', '9':
		return "sh" + symbol, nil
	default:
		return "", fmt.Errorf("unsupported CN stock prefix %q", symbol[:1])
	}
}

func parseTencentTime(value string) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Now().UTC()
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.UTC
	}
	parsed, err := time.ParseInLocation("20060102150405", value, loc)
	if err != nil {
		return time.Now().UTC()
	}
	return parsed.UTC()
}

func scaleEastmoneyPrice(value float64) float64 {
	return value / 100
}
