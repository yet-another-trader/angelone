package angelone

import (
	"context"
	"net/http"
)

// https://smartapi.angelbroking.com/docs/User

type AuthenticateInput struct {
	ClientCode string `json:"clientcode"`
	Password   string `json:"password"`
	TOTP       string `json:"totp"`
	State      string `json:"state"`
}

type AuthenticateOutput struct {
	JwtToken     string `json:"jwtToken"`
	RefreshToken string `json:"refreshToken"`
	FeedToken    string `json:"feedToken"`
	State        string `json:"state"`
}

// Authenticate authenticates the client with the given input.
// Suggested use is to not call this function at all and let the NewClient function do it for you.
func (c *Client) Authenticate(ctx context.Context, input AuthenticateInput) (*Response[AuthenticateOutput], error) {
	method := http.MethodPost
	path := "/auth/angelbroking/user/v1/loginByPassword"

	var res Response[AuthenticateOutput]

	err := c.httpDo(ctx, method, path, input, tradeKey, &res)
	if err != nil {
		return nil, err
	}

	if !res.Status {
		return nil, res
	}

	c.AccessToken = res.Data.JwtToken
	c.RefreshToken = res.Data.RefreshToken
	c.FeedToken = res.Data.FeedToken

	if c.cachePath != nil {
		err = c.cache()
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

type ProfileOutput struct {
	ClientCode    string   `json:"clientcode"`
	Name          string   `json:"name"`
	Email         string   `json:"email"`
	MobileNo      string   `json:"mobileno"`
	Exchanges     []string `json:"exchanges"`
	Products      []string `json:"products"`
	LastLoginTime string   `json:"lastlogintime"`
	Brokerid      string   `json:"brokerid"`
}

// Profile returns the profile of the currently authenticated client.
func (c *Client) Profile(ctx context.Context) (*Response[ProfileOutput], error) {
	method := http.MethodGet
	path := "/secure/angelbroking/user/v1/getProfile"

	var res Response[ProfileOutput]

	err := c.httpDo(ctx, method, path, nil, tradeKey, &res)
	if err != nil {
		return nil, err
	}

	if !res.Status {
		return nil, res
	}

	return &res, nil
}

type RefreshOutput struct {
	JwtToken     string `json:"jwtToken"`
	RefreshToken string `json:"refreshToken"`
	FeedToken    string `json:"feedToken"`
}

// Refresh refreshes the client's access token.
// Suggested use is to not call this function at all and let the NewClient function do it for you.
func (c *Client) Refresh(ctx context.Context) (*Response[RefreshOutput], error) {
	method := http.MethodPost
	path := "/auth/angelbroking/jwt/v1/generateTokens"
	body := struct {
		RefreshToken string `json:"refreshToken"`
	}{RefreshToken: c.RefreshToken}

	var res Response[RefreshOutput]

	err := c.httpDo(ctx, method, path, body, tradeKey, &res)
	if err != nil {
		return nil, err
	}

	if !res.Status {
		return nil, res
	}

	c.AccessToken = res.Data.JwtToken
	c.RefreshToken = res.Data.RefreshToken
	c.FeedToken = res.Data.FeedToken

	if c.cachePath != nil {
		err = c.cache()
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}
