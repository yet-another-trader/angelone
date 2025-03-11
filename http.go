package angelone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var localIP, publicIP, mac string

var client http.Client = http.Client{}

func defaultHeaders(apiKey, accessToken string) (map[string]string, error) {
	var err error

	if localIP == "" {
		localIP, err = localIPAddr()
		if err != nil {
			return nil, err
		}
	}

	if publicIP == "" {
		publicIP, err = publicIPAddr()
		if err != nil {
			return nil, err
		}
	}

	if mac == "" {
		mac, err = macAddr()
		if err != nil {
			return nil, err
		}
	}

	return map[string]string{
		"Content-Type":     "application/json",
		"X-ClientLocalIP":  localIP,
		"X-ClientPublicIP": publicIP,
		"X-MACAddress":     mac,
		"Accept":           "application/json",
		"X-PrivateKey":     apiKey,
		"X-UserType":       "USER",
		"X-SourceID":       "WEB",
		"Authorization":    fmt.Sprintf("Bearer %s", accessToken),
	}, nil
}

type keyType int

const (
	tradeKey keyType = iota
	historyKey
)

func (c *Client) httpDo(ctx context.Context, method, path string, body any, keyType keyType, out any) error {
	j, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", c.basePath, path), bytes.NewBuffer(j))
	if err != nil {
		return err
	}

	apiKey := ""
	switch keyType {
	case tradeKey:
		apiKey = c.TradeKey
	case historyKey:
		apiKey = c.HistoryKey
	}

	headers, err := defaultHeaders(apiKey, c.AccessToken)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, out)
	if err != nil {
		return err
	}

	return nil
}

type Response[T any] struct {
	Status    bool   `json:"status"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorcode"`
	Data      T      `json:"data"`
}

func (r Response[T]) Error() string {
	return fmt.Sprintf("status: %v, message: %s, errorcode: %s", r.Status, r.Message, r.ErrorCode)
}
