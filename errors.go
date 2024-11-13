package goFinance

import (
	"errors"
)

var errTickerNotFound = errors.New("ticker not found")

var errInvalidDateRange = errors.New("invalid date range. Must be 1d, 5d, 1mo, 3mo, 6mo, 1y, 5y, 10y, or ytd")

var errInvalidInterval = errors.New("invalid interval. Must be 1m, 2m, 3m, 5m, 15m, 30m, 60m, 4h, 1d, 1wk, 1mo, 1y")

var errInvalidWindow = errors.New("invalid rolling window. must be greater than 1 or less than length of target column")
