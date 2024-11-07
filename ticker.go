package goFinance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Ticker struct {
	Dates  []int64
	Close  []float64
	Volume []int64
}

type Result struct {
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
}

type Response struct {
	Chart struct {
		Result []Result `json:"result"`
		Error  any
	} `json="chart"`
}

func GetTicker(symbol string) Ticker {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	res := &Response{}
	err = json.Unmarshal(body, res)
	if err != nil {
		panic(err)
	}

	return Ticker{
		Dates:  res.Chart.Result[0].Timestamp,
		Close:  res.Chart.Result[0].Indicators.Quote[0].Close,
		Volume: res.Chart.Result[0].Indicators.Quote[0].Volume,
	}
}
