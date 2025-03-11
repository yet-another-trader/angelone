package angelone

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

// https://smartapi.angelbroking.com/docs/Historical

type Exchange string

const (
	ExchangeNSE Exchange = "NSE" // NSE Stocks and Indices
	ExchangeNFO Exchange = "NFO" // NSE Futures and Options
	ExchangeBSE Exchange = "BSE" // BSE Stocks and Indices
	ExchangeBFO Exchange = "BFO" // BSE Future and Options
	ExchangeCDS Exchange = "CDS" // Currency Derivatives
	ExchangeMCX Exchange = "MCX" // MCX Commodities
)

type Interval string

const (
	Interval1Minute  Interval = "ONE_MINUTE"
	Interval3Minute  Interval = "THREE_MINUTE"
	Interval5Minute  Interval = "FIVE_MINUTE"
	Interval10Minute Interval = "TEN_MINUTE"
	Interval15Minute Interval = "FIFTEEN_MINUTE"
	Interval30Minute Interval = "THIRTY_MINUTE"
	Interval1Hour    Interval = "ONE_HOUR"
	Interval1Day     Interval = "ONE_DAY"
)

var maxDays = map[Interval]int{
	Interval1Minute:  30,
	Interval3Minute:  60,
	Interval5Minute:  100,
	Interval10Minute: 100,
	Interval15Minute: 200,
	Interval30Minute: 200,
	Interval1Hour:    400,
	Interval1Day:     2000,
}

type CandleInput struct {
	Exchange    Exchange  `json:"exchange"`
	SymbolToken string    `json:"symboltoken"`
	Interval    Interval  `json:"interval"`
	FromTime    time.Time `json:"-"`
	ToTime      time.Time `json:"-"`
}

type Candle struct {
	Time   time.Time
	Open   decimal.Decimal
	High   decimal.Decimal
	Low    decimal.Decimal
	Close  decimal.Decimal
	Volume uint
}

type Candles []*Candle

type CandleOutput struct {
	Status    bool   `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorcode"`
	RawData   []any  `json:"data"`
}

func (co CandleOutput) Error() string {
	return fmt.Sprintf("status: %v, message: %s, errorcode: %s", co.Status, co.Message, co.ErrorCode)
}

// Candle fetches historical candle data for a given symbol.
// The maximum number of days that can be fetched is limited by the interval.
// The maximum number of days for each interval is as follows:
//
//	1 minute: 30 days
//	3 minute: 60 days
//	5 minute: 100 days
//	10 minute: 100 days
//	15 minute: 200 days
//	30 minute: 200 days
//	1 hour: 400 days
//	1 day: 2000 days
func (c *Client) Candle(ctx context.Context, input CandleInput) (Candles, error) {
	if input.ToTime.Sub(input.FromTime).Hours()/24 > float64(maxDays[input.Interval]) {
		return nil, errors.New("maximum days exceeded")
	}

	method := http.MethodPost
	path := "/secure/angelbroking/historical/v1/getCandleData"
	body := struct {
		Exchange    Exchange `json:"exchange"`
		SymbolToken string   `json:"symboltoken"`
		Interval    Interval `json:"interval"`
		FromDate    string   `json:"fromdate"`
		ToDate      string   `json:"todate"`
	}{
		Exchange:    input.Exchange,
		SymbolToken: input.SymbolToken,
		Interval:    input.Interval,
		FromDate:    input.FromTime.Format("2006-01-02 15:04"),
		ToDate:      input.ToTime.Format("2006-01-02 15:04"),
	}

	var res CandleOutput

	err := c.httpDo(ctx, method, path, body, tradeKey, &res)
	if err != nil {
		return nil, err
	}

	if !res.Status {
		return nil, res
	}

	layout := "2006-01-02T15:04:05-07:00"
	prices := make([]*Candle, len(res.RawData))
	for i, data := range res.RawData {
		arr, ok := data.([]any)
		if !ok {
			continue
		}

		parsedTime, _ := time.Parse(layout, arr[0].(string))

		prices[i] = &Candle{
			Time:   parsedTime,
			Open:   decimal.NewFromFloat(arr[1].(float64)),
			High:   decimal.NewFromFloat(arr[2].(float64)),
			Low:    decimal.NewFromFloat(arr[3].(float64)),
			Close:  decimal.NewFromFloat(arr[4].(float64)),
			Volume: uint(arr[5].(float64)),
		}
	}

	return prices, nil
}
