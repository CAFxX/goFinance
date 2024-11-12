package goFinance

import (
	"errors"
)

var errTickerNotFound = errors.New("ticker not found")

var errInvalidDateRange = errors.New("invalid date range. Must be 1d, 5d, 1mo, 3mo, 6mo, 1y, 5y, 10y, or ytd")
