# Angel One

A Go library to access AngelOne's SmartAPI.

#### Usage:

```
go get github.com/yet-another-trader/angelone
```

#### Disclaimer:

This package currently **only implements authentication and candle APIs** for Angel One's SmartAPI. Other APIs and functionalities are **not supported** at this time. Use this package with the understanding that its scope is limited, and additional features may need to be implemented separately.

This project is not affiliated with or endorsed by Angel One. Use at your own risk.

#### Example:

```go
package main

import (
    "context"
    "fmt"
    "time"

    angelone "github.com/yet-another-trader/angelone"
)

func main() {
    // required for first login
    clientCode := "<CODE>"
    pin := "<PIN>"
    totp := "<TOTP>"

    // get from https://smartapi.angelbroking.com/apps or https://smartapi.angelbroking.com/create
    tradeKey := "<TRADE-API-KEY>"
    historyKey := "<HISTORY-API-KEY>"

    c := angelone.NewClient(
        context.Background(),
        angelone.WithCachePath("token.json"),
        angelone.WithTradeKey(tradeKey),
        angelone.WithHistoryKey(historyKey),
        angelone.WithAuthenticationInput(angelone.AuthenticateInput{
            ClientCode: clientCode,
            Password:   pin,
            TOTP:       totp,
            State:      "login",
        }),
    )

    co, err := c.Candle(context.Background(), angelone.CandleInput{
        Exchange:     angelone.ExchangeNSE,
        SymbolToken:  "6818",
        Interval:     angelone.Interval1Minute,
        FromTime:     time.Date(2024, 2, 15, 0, 0, 0, 0, time.Local),
        ToTime:       time.Date(2024, 3, 5, 0, 0, 0, 0, time.Local),
    })

    if err != nil {
        panic(err)
    }

    for _, ohlcv := range co {
        fmt.Println(*ohlcv)
    }

    // ...
    // {2024-03-04 15:22:00 +0530 IST 517 517 516.5 516.5 1749}
    // {2024-03-04 15:23:00 +0530 IST 516.45 517 516.3 516.85 2338}
    // {2024-03-04 15:24:00 +0530 IST 516.95 517 516.95 517 1990}
    // {2024-03-04 15:25:00 +0530 IST 516.95 517 516.8 516.85 1039}
    // {2024-03-04 15:26:00 +0530 IST 516.95 517.2 516.85 517.15 1917}
    // {2024-03-04 15:27:00 +0530 IST 517.2 517.2 515.35 515.35 2970}
    // {2024-03-04 15:28:00 +0530 IST 516.5 516.9 515.05 516.9 2209}
    // {2024-03-04 15:29:00 +0530 IST 516.9 517 516.65 516.65 1200}
}
```
