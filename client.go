package angelone

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
)

type Client struct {
	authData     AuthenticateInput
	basePath     string
	cachePath    *string
	ClientCode   string `json:"client_code"`
	TradeKey     string `json:"trade_key,omitempty"`
	HistoryKey   string `json:"history_key,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	FeedToken    string `json:"feed_token,omitempty"`
}

type Option func(*Client)

// WithTradeKey sets the trade key for the client.
// This key is needed for trading operations.
func WithTradeKey(tradeKey string) Option {
	return func(c *Client) {
		c.TradeKey = tradeKey
	}
}

// WithHistoryKey sets the history key for the client.
// This key is needed for historical data operations.
func WithHistoryKey(historyKey string) Option {
	return func(c *Client) {
		c.HistoryKey = historyKey
	}
}

// WithCachePath sets the path to the file where the client will cache its tokens.
// This is useful for avoiding re-authentication on every run.
func WithCachePath(cachePath string) Option {
	return func(c *Client) {
		c.cachePath = &cachePath
	}
}

// WithAuthenticationInput sets the authentication input for the client.
// This is needed for authenticating the client. Ignore this if you have already authenticated and cached your tokens.
func WithAuthenticationInput(input AuthenticateInput) Option {
	return func(c *Client) {
		c.authData = input
	}
}

// NewClient creates a new client with the given options.
// The client will try to authenticate itself using the given authentication input.
func NewClient(ctx context.Context, opts ...Option) *Client {
	c := &Client{
		basePath: "https://apiconnect.angelone.in/rest",
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.cachePath != nil {
		_ = c.load()
	}

	// try getting profile
	_, err := c.Profile(ctx)
	if err == nil {
		// all good with current token, proceed without authentication or refreshing
		return c
	}

	// try refreshing token
	_, err = c.Refresh(ctx)
	if err == nil {
		// all good with refreshed token, proceed without authentication
		return c
	}

	// try authenticating
	_, err = c.Authenticate(ctx, c.authData)
	if err != nil {
		// nothing worked
		panic(err)
	}

	return c
}

func (c *Client) cache() error {
	if c.cachePath == nil {
		return errors.New("cache path not set")
	}

	f, err := os.OpenFile(*c.cachePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) load() error {
	if c.cachePath == nil {
		return errors.New("cache path not set")
	}

	f, err := os.Open(*c.cachePath)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, c)
	if err != nil {
		return err
	}

	return nil
}
