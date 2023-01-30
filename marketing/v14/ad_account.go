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

// FB Business Account represents an ad account.
type FBBusinessAccount struct {
	Id   string `json:"id"`
	Name string `json:"name"`
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
