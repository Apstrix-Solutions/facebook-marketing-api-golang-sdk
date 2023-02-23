package v14

import (
	"context"
	"fmt"

	"github.com/Apstrix-Solutions/facebook-marketing-api-golang-sdk/fb"
)

// AdAccountService works with ad accounts.
type AdAccountService struct {
	c *fb.Client
}

func (aas *AdAccountService) FBBusinessAccount(ctx context.Context) ([]FBBusinessAccount, error) {
	res := []FBBusinessAccount{}
	rb := fb.NewRoute(Version, "/me/businesses").Fields("id", "name")
	fmt.Println(rb)
	err := aas.c.GetList(ctx, rb.String(), &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// List lists all Facebook ad accounts that belong to this business.
func (aas *AdAccountService) List(ctx context.Context, businessID string) ([]AdAccount, error) {
	res := []AdAccount{}
	rb := fb.NewRoute(Version, "/%s/owned_ad_accounts", businessID).Limit(1000).Fields("name", "currency", "account_id")
	err := aas.c.GetList(ctx, rb.String(), &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// List lists all Instagram ad accounts that belong to this business.
func (aas *AdAccountService) ListInstaAccounts(ctx context.Context, businessID string) ([]InstaAdAccount, error) {
	res := []InstaAdAccount{}

	// rb := fb.NewRoute(Version, "/%s/owned_instagram_accounts", businessID).Limit(1000).Fields("id", "username")
	rb := fb.NewRoute(Version, "/%s/instagram_accounts", businessID).Limit(1000).Fields("id", "username")

	err := aas.c.GetList(ctx, rb.String(), &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// List lists all Instagram ad accounts that belong to this business.
func (aas *AdAccountService) LongLivedUserToken(ctx context.Context, req *ShortTokenInfo) (*LongTokenInfo, error) {

	resp := LongTokenInfo{}

	err := aas.c.PostJSON(ctx, fb.NewRoute(Version, "/oauth/access_token").String(), req, &resp)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

type ShortTokenInfo struct {
	GrantType    string `json:"grant_type,omitempty"`
	AccessToken  string `json:"fb_exchange_token,omitempty"` //User Access Token
	ClientId     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

type LongTokenInfo struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// FB Business Account represents an ad account.
type FBBusinessAccount struct {
	Id   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// AdAccount represents an Facebook ad account.
type AdAccount struct {
	Name      string `json:"name"`
	AccountID string `json:"account_id"`
	Currency  string `json:"currency"`
}

// InstaAccount represents an Instagram ad account.
type InstaAdAccount struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}
