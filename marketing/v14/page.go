package v14

import (
	"context"
	"fmt"

	"github.com/Apstrix-Solutions/facebook-marketing-api-golang-sdk/fb"
)

// PageService contains all methods for working on pages.
type PageService struct {
	c *fb.Client
}

// SetPageAccessToken tries to retrieve the access token for a facebook page and includes it in the passed context so the fb.Client can use it for making requests.
func (ps *PageService) SetPageAccessToken(ctx context.Context, pageID string) (context.Context, error) {
	tc := struct {
		AccessToken string `json:"access_token"`
	}{}

	fmt.Println("Url is : ", fb.NewRoute(Version, "/%s", pageID).Fields("access_token").String())

	err := ps.c.GetJSON(ctx, fb.NewRoute(Version, "/%s", pageID).Fields("access_token").String(), &tc)

	if err != nil {
		return ctx, err
	} else if tc.AccessToken == "" {
		return ctx, fmt.Errorf("could not get page access token for '%s'", pageID)
	}

	return fb.SetPageAccessToken(ctx, tc.AccessToken), nil
}

// GetPageBackedInstagramAccounts returns the instagram actor associated with a facebook page.
func (ps *PageService) GetPageBackedInstagramAccounts(ctx context.Context, pageID string) (*InstagramActor, error) {
	ctx, err := ps.SetPageAccessToken(ctx, pageID)
	if err != nil {
		return nil, err
	}

	fpiga := struct {
		PageBackedInstagramAccounts struct {
			Data []InstagramActor `json:"data"`
		} `json:"page_backed_instagram_accounts"`
	}{}
	err = ps.c.GetJSON(ctx, fb.NewRoute(Version, "/%s", pageID).Fields("page_backed_instagram_accounts{id,username}").String(), &fpiga)
	if err != nil {
		return nil, err
	}

	if len(fpiga.PageBackedInstagramAccounts.Data) != 1 {
		return nil, fmt.Errorf("could not get consistent page_backed_instagram_accounts data for facebook page with external id %s", pageID)
	}

	res := fpiga.PageBackedInstagramAccounts.Data[0]
	if res.ID == "" {
		return nil, fmt.Errorf("could not get page_backed_instagram_accounts ID for facebook page with external id %s", pageID)
	}
	if res.Username == "" {
		return nil, fmt.Errorf("could not get page_backed_instagram_accounts username for facebook page with external id %s", pageID)
	}

	return &res, nil
}

// GetClientPages returns all client pages.
func (ps *PageService) GetClientPages(ctx context.Context, businessID string) ([]Page, error) {
	res := []Page{}

	fmt.Println(fb.NewRoute(Version, "/%s/client_pages", businessID))
	route := fb.NewRoute(Version, "/%s/client_pages", businessID).Limit(1000).Fields(pageFields...)
	err := ps.c.GetList(ctx, route.String(), &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetOwnedPages returns all owned pages.
func (ps *PageService) GetOwnedPages(ctx context.Context, businessID string) ([]Page, error) {
	res := []Page{}
	route := fb.NewRoute(Version, "/%s/owned_pages", businessID).Limit(1000).Fields(pageFields...)
	err := ps.c.GetList(ctx, route.String(), &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetInstagramActors returns all instagram accounts.
func (ps *PageService) GetInstagramActors(ctx context.Context, businessID string) ([]InstagramActor, error) {
	res := []InstagramActor{}
	route := fb.NewRoute(Version, "/%s/instagram_accounts", businessID).Limit(1000).Fields(instagramActorFields...)
	err := ps.c.GetList(ctx, route.String(), &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Get returns a single page.
func (ps *PageService) Get(ctx context.Context, id string) (*Page, error) {
	res := &Page{}
	route := fb.NewRoute(Version, "/%s", id).Fields(pageFields...)
	err := ps.c.GetJSON(ctx, route.String(), res)
	if err != nil {
		if fb.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return res, nil
}

// GetInstagramActor returns a single instagram actor.
func (ps *PageService) GetInstagramActor(ctx context.Context, id string) (*InstagramActor, error) {
	res := &InstagramActor{}
	route := fb.NewRoute(Version, "/%s", id).Fields(instagramActorFields...)
	err := ps.c.GetJSON(ctx, route.String(), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

var (
	pageFields           = []string{"id", "global_brand_page_name"}
	instagramActorFields = []string{"id", "username"}
)

// Page represents a facebook page.
type Page struct {
	ID                  string `json:"id"`
	GlobalBrandPageName string `json:"global_brand_page_name"`
}

// InstagramActor represents an instagram actor.
type InstagramActor struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// Lead Generation Form
type LeadGen struct {
	Name              string `json:"name"`
	PageID            string `json:"page_id"`
	FollowUpActionURL string `json:"follow_up_action_url"`
	PrivacyPolicy     struct {
		URL      string `json:"url"`
		LinkText string `json:"link_text"`
	} `json:"privacy_policy"`
	ContextCard struct {
		Title      string   `json:"title"`
		Content    []string `json:"content"`
		Style      string   `json:"style"`
		CoverPhoto string   `json:"cover_photo_id"`
	} `json:"context_card"`
	ThankYouPage struct {
		Title      string `json:"title"`
		Body       string `json:"body"`
		ButtonType string `json:"button_type"`
		ButtonText string `json:"button_text"`
		WebsiteURL string `json:"website_url"`
	} `json:"thank_you_page"`
	Questions []struct {
		Key   string `json:"key"`
		Type  string `json:"type"`
		Label string `json:"label,omitempty"`
	} `json:"questions"`
}

type LeadFormResp struct {
	FormId string `json:"id"`
}

// Create LeadGen Form
func (ps *PageService) CreateLeadAdForm(ctx context.Context, content LeadGen, id string) (*LeadFormResp, error) {

	req := content
	resp := LeadFormResp{}
	err := ps.c.PostJSON(ctx, fb.NewRoute(Version, "/%s", id+"/leadgen_forms").String(), req, &resp)

	if err != nil {
		fmt.Println("Post Json Error : ", err)
		return nil, err
	}

	return &resp, nil

}
