package goFinance

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"slices"
	"time"

	"github.com/corpix/uarand"
)

type Ticker struct {
	Dates      []time.Time
	Indicators map[string][]float64
}

type Response struct {
	Chart struct {
		Result []struct {
			Meta       any     `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Close  []float64 `json:"close"`
					Low    []float64 `json:"low"`
					High   []float64 `json:"high"`
					Open   []float64 `json:"open"`
					Volume []float64 `json:"volume"`
				} `json:"quote"`
				AdjClose []struct {
					AdjClose []float64 `json:"adjclose"`
				} `json:"adjclose"`
			} `json:"indicators"`
		} `json:"result"`
		Error any
	} `json:"chart"`
}

func GetTicker(symbol string, period string, interval string) (Ticker, error) {
	start, end, err := dateRange(period)
	if err != nil {
		return Ticker{}, err
	}

	err = validInterval(interval)
	if err != nil {
		return Ticker{}, err
	}

	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?period1=%d&period2=%d&interval=%s", symbol, start, end, interval)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Ticker{}, err
	}
	req.Header.Set("User-Agent", uarand.GetRandom())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Ticker{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Ticker{}, err
	}

	res := &Response{}

	err = json.Unmarshal(body, res)
	if err != nil {
		return Ticker{}, err
	}

	if len(res.Chart.Result) == 0 {
		return Ticker{}, errTickerNotFound
	}

	data := res.Chart.Result[0]

	ticker := Ticker{
		Dates: parseDates(data.Timestamp),
		Indicators: map[string][]float64{
			"open":   data.Indicators.Quote[0].Open,
			"high":   data.Indicators.Quote[0].High,
			"low":    data.Indicators.Quote[0].Low,
			"close":  data.Indicators.Quote[0].Close,
			"volume": data.Indicators.Quote[0].Volume,
		},
	}

	if len(data.Indicators.AdjClose) > 0 {
		ticker.Indicators["adjClose"] = data.Indicators.AdjClose[0].AdjClose
	}

	return ticker, nil
}

func parseDates(unixTimes []int64) []time.Time {
	res := make([]time.Time, len(unixTimes))

	for i, t := range unixTimes {
		res[i] = time.Unix(t, 0)
	}

	return res
}

func dateRange(period string) (int64, int64, error) {
	var start time.Time
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	switch period {
	case "1d":
		start = today.AddDate(0, 0, -1)
	case "5d":
		start = today.AddDate(0, 0, -5)
	case "1mo":
		start = today.AddDate(0, -1, 0)
	case "3mo":
		start = today.AddDate(0, -3, 0)
	case "6mo":
		start = today.AddDate(0, -6, 0)
	case "1y":
		start = today.AddDate(-1, 0, 0)
	case "2y":
		start = today.AddDate(-2, 0, -5)
	case "5y":
		start = today.AddDate(-5, 0, -5)
	case "10y":
		start = today.AddDate(-10, 0, -5)
	case "ytd":
		start = time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	case "max":
		start = time.Unix(math.MinInt32, 0)
	default:
		// default to 1y prior
		return 0, 0, errInvalidDateRange
	}
	return start.Unix(), today.Unix(), nil
}

func validInterval(interval string) error {
	intervals := []string{
		"1m",
		"2m",
		"3m",
		"5m",
		"15m",
		"30m",
		"60m",
		"4h",
		"1d",
		"1wk",
		"1mo",
		"1y",
	}

	if !slices.Contains(intervals, interval) {
		return errInvalidInterval
	}

	return nil
}

func (t *Ticker) RollingAverage(target string, result string, window int) error {
	if window < 1 || window >= len(t.Indicators[target]) {
		return errInvalidWindow
	}
	vals := t.Indicators[target]
	newCol := make([]float64, len(vals))

	for idx, _ := range vals {
		if idx < window-1 {
			newCol[idx] = 0
		} else {
			sum := 0.0
			slice := vals[idx-window+1 : idx+1]
			for _, x := range slice {
				sum += x
			}

			newCol[idx] = sum / float64(len(slice))
		}
	}

	t.Indicators[result] = newCol

	return nil
}
