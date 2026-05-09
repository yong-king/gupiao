package marketdata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var tencentKLineURL = "https://web.ifzq.gtimg.cn/appstock/app/fqkline/get"

type tencentKLineResponse struct {
	Code int `json:"code"`
	Data map[string]struct {
		QFQDay [][]string `json:"qfqday"`
		Day    [][]string `json:"day"`
	} `json:"data"`
}

func FetchTencentDailyKLines(ctx context.Context, request QuoteRequest, limit int, client *http.Client) ([]KLine, error) {
	if limit <= 0 {
		limit = 60
	}
	if limit > 240 {
		limit = 240
	}
	symbol := strings.ToUpper(strings.TrimSpace(request.Symbol))
	code, err := tencentSymbol(symbol)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(tencentKLineURL)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("param", fmt.Sprintf("%s,day,,,%d,qfq", code, limit))
	u.RawQuery = query.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Referer", "https://gu.qq.com/")
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 jijin-stock-monitor/0.1")
	if client == nil {
		client = &http.Client{Timeout: 8 * time.Second}
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("tencent kline returned status %d", resp.StatusCode)
	}

	var parsed tencentKLineResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}
	if parsed.Code != 0 {
		return nil, fmt.Errorf("tencent kline returned code %d", parsed.Code)
	}
	rows := parsed.Data[code].QFQDay
	if len(rows) == 0 {
		rows = parsed.Data[code].Day
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("tencent kline missing rows for %s", request.Key())
	}
	out := make([]KLine, 0, len(rows))
	for _, row := range rows {
		if len(row) < 6 {
			continue
		}
		open, _ := parseFloat(row[1])
		closePrice, _ := parseFloat(row[2])
		high, _ := parseFloat(row[3])
		low, _ := parseFloat(row[4])
		volume, _ := strconv.ParseFloat(row[5], 64)
		dataTime := parseKLineDate(row[0])
		out = append(out, KLine{
			Market:   "CN",
			Symbol:   symbol,
			Date:     row[0],
			Open:     open,
			Close:    closePrice,
			High:     high,
			Low:      low,
			Volume:   int64(volume * 100),
			Source:   "tencent-kline",
			DataTime: dataTime,
		})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("tencent kline rows invalid for %s", request.Key())
	}
	return out, nil
}

func parseKLineDate(value string) time.Time {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.UTC
	}
	parsed, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(value), loc)
	if err != nil {
		return time.Now().UTC()
	}
	return parsed.UTC()
}
