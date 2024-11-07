package goFinance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Ticker struct {
	Dates  []time.Time
	Close  []float64
	Volume []int64
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
					Volume []int64   `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error any
	} `json:"chart"`
}

func GetTicker(symbol string) (Ticker, error) {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s", symbol)

	resp, err := http.Get(url)
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

	return Ticker{
		Dates:  parseDates(data.Timestamp),
		Close:  data.Indicators.Quote[0].Close,
		Volume: data.Indicators.Quote[0].Volume,
	}, nil
}

func parseDates(unixTimes []int64) []time.Time {
	res := make([]time.Time, len(unixTimes))

	for i, t := range unixTimes {
		res[i] = time.Unix(t, 0)
	}

	return res
}
