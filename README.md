# goFinance
Go package to access the unofficial Yahoo Finance API

### Overview
This package current produces a 'dataframe' (very loose description) containing dates and indicators:

```go
type Ticker struct {
	Dates      []time.Time
	Indicators map[string][]float64
}
```

The indicators are a string indexed map, which so far includes Open, High, Low, Close, Adjusted Close, and Volume data for a given time period and interval.

The time period and interval come from Yahoo Finance's allowed values, so these are things like 1d, 1mo, 1y etc. Keep in mind this is all free data from Yahoo so while you can get minute-by-minute data I would not run a bot off this (it's delayed at least 15 minutes)

There's one helper method added to the Ticker struct to allow for calculating a simple rolling average. It takes the target data name, the new column's name, and the window to look back on. I'll be adding additional 'indicator' functions over time, but this is it so far.




### Usage

`go get github.com/joetats/goFinance`

```go
package main

import (
	"fmt"
	"github.com/joetats/goFinance"
)

func main() {
	df, err := goFinance.GetTicker("AAPL", "1mo, "1d)
	if err != nil {
		panic(err)
    }
	
	err = df.RollingAverage("close", "rollingAverageClose", 7)
	if err != nil {
		panic(err)
    }
	
	fmt.Printf(`%v`, df)
}
```