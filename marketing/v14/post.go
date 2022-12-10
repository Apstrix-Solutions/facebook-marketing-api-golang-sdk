package v14

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/justwatchcom/facebook-marketing-api-golang-sdk/fb"
)

// PostService works on posts.
type PostService struct {
	c *fb.Client
	*fb.StatsContainer
}

// ------------------------------------------------ Facebook Models Starts ------------------------------------------------ //
type RespPost struct {
	Posts struct {
		Data []struct {
			ID          string      `json:"id"`
			Message     string      `json:"message"`
			CreatedTime string      `json:"created_time"`
			Attachments Attachments `json:"attachments"`
		} `json:"data"`
	} `json:"posts"`
}

type RespPost_Paid struct {
	AdPosts struct {
		Data []struct {
			ID          string      `json:"id"`
			Message     string      `json:"message"`
			CreatedTime string      `json:"created_time"`
			Attachments Attachments `json:"attachments"`
		} `json:"data"`
	} `json:"ads_posts"`
}

type Attachments struct {
	Data []AttachmentData `json:"data"`
}

type AttachmentData struct {
	Description string `json:"description"`
	Media       struct {
		Image struct {
			Height int    `json:"height"`
			Src    string `json:"src"`
			Width  int    `json:"width"`
		} `json:"image"`
	} `json:"media"`
	Target struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	} `json:"target"`
	Type string `json:"type"`
	URL  string `json:"url,omitempty"`
}

type GetPostContent struct {
	Message string   `json:"message"`
	Urls    []string `json:"url"`
}

type PostPhotoContent struct {
	PhotoUrl  string `json:"url"`
	Published bool   `json:""`
}

type PostPhotoAttachments struct {
	Message     string   `json:"message"`
	Attachments []string `json:"attached_media"`
}

type PostPhotoReponse struct {
	Photo_ID string `json:"id"`
}

type PostMsgReponse struct {
	Post_ID string `json:"id"`
}

type PostVideoContent struct {
	FileUrl     string `json:"file_url"`
	Description string `json:"description"`
}

// ------------------------------------------------ Facebook Models Ends------------------------------------------------ //

// ------------------------------------------------ Instagram Models Starts------------------------------------------------ //

type InstaPostSingleContent struct {
	MediaType    string `json:"media_type,omitempty"`
	Caption      string `json:"caption,omitempty"`
	VideoUrl     string `json:"video_url,omitempty"`
	ImageUrl     string `json:"image_url,omitempty"`
	Timestamp    string `json:"timestamp,omitempty"`
	Published    bool   `json:"published,omitempty"`
	ScheduleDate string `json:"scheduled_publish_time,omitempty"`
}

type InstaPostReponse struct {
	Post_ID string `json:"id"`
}

type ContainerData struct {
	ID []string `json:"id"`
}

// ------------------------------------------------ Instagram Models Ends------------------------------------------------ //

// Instagram Publish Single Media Post
func (ps *PostService) Create_SinglePostContainer(ctx context.Context, instaId string, content *InstaPostSingleContent) (*InstaPostReponse, error) {

	req := content
	resp := InstaPostReponse{}

	err := ps.c.PostJSON(ctx, fb.NewRoute(Version, "/%s", instaId+"/media").String(), req, &resp)

	if err != nil {
		fmt.Println("Post Json Error : ", err)
		return nil, err
	}

	return &resp, nil
}

func (ps *PostService) GetInstaOrganicPosts(ctx context.Context, instaId string, urlParams []string) (*[]RespPost, error) {

	res := &RespPost{}
	allPosts := []RespPost{}

	err := ps.c.GetJSON(ctx, fb.NewRoute(Version, "/%s", instaId).Fields(urlParams...).String(), res)
	fmt.Println(fb.NewRoute(Version, "/%s", instaId).Fields(urlParams...).String())
	// err := ps.c.GetJSON(ctx, fb.NewRoute(Version, "/%s", businessId).Fields("posts{id, message, attachments}").String(), res)

	if err != nil {
		if fb.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	allPosts = append(allPosts, *res)

	return &allPosts, nil
}

func (ps *PostService) GetInstaPaidPosts(ctx context.Context, instaId string, urlParams []string) (*[]RespPost_Paid, error) {

	res := &RespPost_Paid{}
	allPosts := []RespPost_Paid{}

	err := ps.c.GetJSON(ctx, fb.NewRoute(Version, "/%s", instaId).Fields(urlParams...).String(), res)

	fmt.Println(fb.NewRoute(Version, "/%s", instaId).Fields(urlParams...).String())
	// err := ps.c.GetJSON(ctx, fb.NewRoute(Version, "/%s", businessId).Fields("posts{id, message, attachments}").String(), res)

	if err != nil {
		if fb.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	allPosts = append(allPosts, *res)

	return &allPosts, nil
}

func (ps *PostService) GetInstaPost(ctx context.Context, id string) (*Post, error) {
	res := &Post{}
	err := ps.c.GetJSON(ctx, fb.NewRoute(Version, "/%s", id).Fields(postFields...).String(), res)
	if err != nil {
		if fb.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}
	err = ps.getPostAttachments(ctx, res)
	if err != nil {
		return nil, err
	}
	if res.Type != "video" && res.Type != "photo" && res.Type != "link" {
		res.Type = "unknown"
	}

	return res, nil
}

//Publish Instagram Container/s

// Publish a Content and Return Response with Post ID
func (ps *PostService) CreateMsgPost(ctx context.Context, businessId string, content GetPostContent, pageAccessToken string) (*PostMsgReponse, error) {

	route := fb.NewRoute(Version, "/%s", businessId+"/feed").String()
	data := url.Values{}

	data.Add("message", content.Message)
	data.Add("access_token", pageAccessToken)

	response, err := http.PostForm(route, data)

	if err != nil {

		fmt.Println("Post Json Error : ", err)
		return nil, err
	}

	r := PostMsgReponse{}
	json.NewDecoder(response.Body).Decode(&r)
	defer response.Body.Close()

	return &r, nil
}

// Publish a Content (Photo) and Return Response with Photo ID
func (ps *PostService) UploadPhotoFB(ctx context.Context, businessId string, imageUrl string, pageAccessToken string) (*PostPhotoReponse, error) {

	route := fb.NewRoute(Version, "/%s", businessId+"/photos").String()
	data := url.Values{}

	data.Add("url", imageUrl)
	data.Add("published", "false")
	data.Add("access_token", pageAccessToken)

	response, err := http.PostForm(route, data)

	if err != nil {

		fmt.Println("Post Json Error : ", err)
		return nil, err
	}

	r := PostPhotoReponse{}
	json.NewDecoder(response.Body).Decode(&r)
	defer response.Body.Close()

	return &r, nil
}

// Publish a Content (Photo) and Return Response with Photo ID
func (ps *PostService) PostContentWithPhotos(ctx context.Context, businessId string, content PostPhotoAttachments, pageAccessToken string) (*PostPhotoReponse, error) {

	route := fb.NewRoute(Version, "/%s", businessId+"/feed").String()
	data := url.Values{}

	ids, err := json.Marshal(content.Attachments)

	if err != nil {
		fmt.Println("Marshalling Error : ", err)
		return nil, err
	}

	data.Add("message", content.Message)
	data.Add("attached_media", string(ids))
	data.Add("access_token", pageAccessToken)

	response, err := http.PostForm(route, data)

	if err != nil {
		fmt.Println("Post Json Error : ", err)
		return nil, err
	}

	r := PostPhotoReponse{}

	json.NewDecoder(response.Body).Decode(&r)
	defer response.Body.Close()

	return &r, nil
}

// Publish a Content (Video) and Return Response with Post Id
func (ps *PostService) CreateVideoPost(ctx context.Context, businessId string, content PostVideoContent, pageAccessToken string) (*PostPhotoReponse, error) {

	route := fb.NewRoute(Version, "/%s", businessId+"/videos").String()
	data := url.Values{}

	data.Add("file_url", content.FileUrl)
	data.Add("description", content.Description)
	data.Add("access_token", pageAccessToken)

	response, err := http.PostForm(route, data)

	if err != nil {
		fmt.Println("Post Json Error : ", err)
		return nil, err
	}

	r := PostPhotoReponse{}

	json.NewDecoder(response.Body).Decode(&r)
	defer response.Body.Close()

	return &r, nil
}

// Get returns a single post.
func (ps *PostService) Get(ctx context.Context, id string) (*Post, error) {
	res := &Post{}
	err := ps.c.GetJSON(ctx, fb.NewRoute(Version, "/%s", id).Fields(postFields...).String(), res)
	if err != nil {
		if fb.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}
	err = ps.getPostAttachments(ctx, res)
	if err != nil {
		return nil, err
	}
	if res.Type != "video" && res.Type != "photo" && res.Type != "link" {
		res.Type = "unknown"
	}

	return res, nil
}

// GetPosts returns all posts under an account
func (ps *PostService) GetPosts(ctx context.Context, businessId string, urlParams []string) (*[]RespPost, error) {

	res := &RespPost{}
	allPosts := []RespPost{}

	err := ps.c.GetJSON(ctx, fb.NewRoute(Version, "/%s", businessId).Fields(urlParams...).String(), res)
	// err := ps.c.GetJSON(ctx, fb.NewRoute(Version, "/%s", businessId).Fields("posts{id, message, attachments}").String(), res)

	if err != nil {
		if fb.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	allPosts = append(allPosts, *res)

	return &allPosts, nil
}

func (ps *PostService) getPostAttachments(ctx context.Context, post *Post) error {
	if post == nil {
		return nil
	}
	pAs := &StoryAttachments{}
	err := ps.c.GetJSON(ctx, fb.NewRoute(Version, "/%s/attachments", post.ID).Fields(postAttachmentsFields...).String(), pAs)
	if err != nil {
		if fb.IsNotFound(err) {
			return nil
		}

		return err
	}
	pA := StoryAttachment{}
	if len(pAs.Data) > 0 {
		pA = pAs.Data[0]
	} else {
		return nil
	}
	sat := strings.TrimSpace(pA.StoryAttachmentType)
	mt := strings.TrimSpace(pA.MediaType)
	if mt != "" {
		if mt == "link" {
			post.Type = "status"
		} else {
			post.Type = mt
		}
	} else if sat != "" {
		post.Type = sat
	}
	if strings.HasPrefix(post.Type, "video") {
		post.Type = "video"
	}
	if post.Type == "album" {
		post.Type = "photo"
	}
	if pA.Title != "" {
		post.Name = pA.Title
	} else if pA.Name != "" {
		post.Name = pA.Name
	}
	if pA.URL != "" {
		post.Link = pA.URL
	} else if pA.UnshimmedURL != "" {
		post.Link = pA.UnshimmedURL
	}
	if pA.Description != "" {
		post.Description = pA.Description
	}
	if pA.Target != nil {
		if pA.Target.ID != "" {
			post.ObjectID = pA.Target.ID
		}
	}

	return nil
}

// GetReactions returns the amount of reactions for a post.
func (ps *PostService) GetReactions(ctx context.Context, postID string) (Reactions, error) {
	m := Reactions{}
	for _, r := range reactions {
		rc := fb.SummaryContainer{}
		err := ps.c.GetJSON(ctx, fb.NewRoute(Version, "/%s/reactions", postID).Summary("total_count").Limit(0).Type(r).String(), &rc)
		if err != nil {
			return nil, err
		} else if rc.Summary.TotalCount == 0 {
			continue
		}

		m[r] = rc.Summary.TotalCount
	}

	return m, nil
}

// CountComments returns the total amount of parent comments.
func (ps *PostService) CountComments(ctx context.Context, postID string) (uint64, error) {
	sc := &fb.SummaryContainer{}
	err := ps.c.GetJSON(ctx, fb.NewRoute(Version, "/%s/comments", postID).Limit(0).Summary("1").String(), sc)

	return sc.Summary.TotalCount, err
}

// ListComments creates a new CommentListCall
// Filters may be "stream" or "toplevel".
func (ps *PostService) ListComments(postID, filter string) *CommentListCall {
	return &CommentListCall{
		RouteBuilder:   fb.NewRoute(Version, "/%s/comments", postID).Fields("message", "message_tags", "parent", "from", "created_time").Limit(100).Order("chronological").Filter(filter),
		c:              ps.c,
		id:             postID,
		StatsContainer: ps.StatsContainer,
	}
}

// CommentListCall is used for listing comments of a post.
type CommentListCall struct {
	*fb.RouteBuilder
	c  *fb.Client
	id string
	*fb.StatsContainer
}

// List performs the CommentListCall and returns all comments as slice.
func (clc *CommentListCall) List(ctx context.Context) ([]Comment, error) {
	stats := clc.StatsContainer.AddStats(clc.id)
	if stats == nil {
		return nil, fmt.Errorf("post %s comments already being downloaded", clc.id)
	}
	defer clc.StatsContainer.RemoveStats(clc.id)
	ctx = stats.AddToContext(ctx)
	res := []Comment{}
	err := clc.c.GetList(ctx, clc.RouteBuilder.String(), &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (clc *CommentListCall) Read(ctx context.Context, c chan<- Comment) error {
	stats := clc.StatsContainer.AddStats(clc.id)
	if stats == nil {
		return fmt.Errorf("post %s comments already being downloaded", clc.id)
	}

	jres := make(chan json.RawMessage)
	wg := errgroup.Group{}
	wg.Go(func() error {
		defer close(jres)

		return clc.c.ReadList(ctx, clc.RouteBuilder.String(), jres)
	})
	wg.Go(func() error {
		for e := range jres {
			v := Comment{}
			err := json.Unmarshal(e, &v)
			if err != nil {
				return err
			}
			stats.Add(1)
			c <- v
		}
		clc.StatsContainer.RemoveStats(clc.id)

		return nil
	})

	return wg.Wait()
}

// Other fields that can be used:
// "actions",
// "admin_creator",
// "allowed_advertising_objectives",
// "application",
// "backdated_time",
// "caption",
// "child_attachments",
// "comments_mirroring_domain",
// "coordinates",
// "created_time",
// "event",
// "expanded_height",
// "expanded_width",
// "feed_targeting",
// "full_picture",
// "height",
// "icon",
// "instagram_eligibility",
// "is_app_share",
// "is_expired",
// "is_hidden",
// "is_instagram_eligible",
// "is_popular",
// "is_published",
// "is_spherical",
// "message_tags",
// "multi_share_end_card",
// "multi_share_optimized",
// "parent_id",
// "place",
// "privacy",
// "promotion_status",
// "properties",
// "scheduled_publish_time",
// "shares",
// "status_type",
// "story",
// "story_tags",
// "subscribed",
// "target",
// "targeting",
// "timeline_visibility",
// "updated_time",
// "via",
// "video_buying_eligibility",
// "width",.
var (
	postFields            = []string{"call_to_action", "from", "id", "message", "picture", "promotable_id"}
	reactions             = []string{"LIKE", "LOVE", "WOW", "HAHA", "SAD", "ANGRY", "THANKFUL"}
	postAttachmentsFields = []string{"description", "name", "type", "url", "target", "media_type"}
)

// Post represents the fb graph api response for a fb video post https://developers.facebook.com/docs/graph-api/reference/v5.0/page-post
type Post struct {
	CreatedTime            string                          `json:"created_time"`
	ContentCategory        string                          `json:"content_category"`
	Description            string                          `json:"description"`
	EmbedHTML              string                          `json:"embed_html"`
	Embeddable             bool                            `json:"embeddable"`
	ID                     string                          `json:"id"`
	Icon                   string                          `json:"icon"`
	IsInstagramEligible    bool                            `json:"is_instagram_eligible"`
	Picture                string                          `json:"picture"`
	PermalinkURL           string                          `json:"permalink_url"`
	MonetizationStatus     string                          `json:"monetization_status"`
	Length                 float64                         `json:"length"`
	Link                   string                          `json:"link"`
	Name                   string                          `json:"name"`
	Type                   string                          `json:"type"`
	Published              bool                            `json:"published"`
	UpdatedTime            string                          `json:"updated_time"`
	Message                string                          `json:"message"`
	InstagramEligibility   string                          `json:"instagram_eligibility"`
	FullPicture            string                          `json:"full_picture"`
	MultiShareEndCard      bool                            `json:"multi_share_end_card"`
	MultiShareOptimized    bool                            `json:"multi_share_optimized"`
	ObjectID               string                          `json:"object_id"`
	PromotableID           string                          `json:"promotable_id"`
	PromotionStatus        string                          `json:"promotion_status"`
	StatusType             string                          `json:"status_type"`
	Subscribed             bool                            `json:"subscribed"`
	TimelineVisibility     string                          `json:"timeline_visibility"`
	VideoBuyingEligibility []string                        `json:"video_buying_eligibility"`
	IsHidden               bool                            `json:"is_hidden"`
	IsAppShare             bool                            `json:"is_app_share"`
	IsExpired              bool                            `json:"is_expired"`
	IsPopular              bool                            `json:"is_popular"`
	IsPublished            bool                            `json:"is_published"`
	IsSpherical            bool                            `json:"is_spherical"`
	CallToAction           *AdCreativeLinkDataCallToAction `json:"call_to_action"`
	Format                 []struct {
		EmbedHTML string `json:"embed_html"`
		Filter    string `json:"filter"`
		Height    int    `json:"height"`
		Picture   string `json:"picture"`
		Width     int    `json:"width"`
	} `json:"format"`
	From    IDContainer `json:"from"`
	Privacy struct {
		Allow       string `json:"allow"`
		Deny        string `json:"deny"`
		Description string `json:"description"`
		Friends     string `json:"friends"`
		Networks    string `json:"networks"`
		Value       string `json:"value"`
	} `json:"privacy"`
	Status struct {
		VideoStatus string `json:"video_status"`
	} `json:"status"`
	Application struct {
		Category string `json:"category"`
		Link     string `json:"link"`
		Name     string `json:"name"`
		ID       string `json:"id"`
	} `json:"application"`
	Coordinates struct{} `json:"coordinates"`
	Actions     []struct {
		Name string `json:"name"`
		Link string `json:"link"`
	} `json:"actions"`
	Properties []struct {
		Name string `json:"name"`
		Text string `json:"text"`
	} `json:"properties"`
}

// StoryAttachment holds information about a post, used since v3.3 https://developers.facebook.com/docs/graph-api/reference/story-attachment/
type StoryAttachment struct {
	Description         string                 `json:"description,omitempty"`
	Media               *StoryAttachmentMedia  `json:"media,omitempty"`
	MediaType           string                 `json:"media_type,omitempty"`
	Title               string                 `json:"title,omitempty"`
	StoryAttachmentType string                 `json:"type,omitempty"`
	UnshimmedURL        string                 `json:"unshimmed_url,omitempty"`
	URL                 string                 `json:"url,omitempty"`
	Name                string                 `json:"name,omitempty"`
	Target              *StoryAttachmentTarget `json:"target,omitempty"`
}

// StoryAttachments wraps the data slice around the StoryAttachment(s).
type StoryAttachments struct {
	Data []StoryAttachment `json:"data,omitempty"`
}

// StoryAttachmentTarget https://developers.facebook.com/docs/graph-api/reference/story-attachment-target/
type StoryAttachmentTarget struct {
	ID           string `json:"id,omitempty"`
	UnshimmedURL string `json:"unshimmed_url,omitempty"`
	URL          string `json:"url,omitempty"`
}

// StoryAttachmentMedia https://developers.facebook.com/docs/graph-api/reference/v5.0/story-attachment-media
type StoryAttachmentMedia struct {
	Image  interface{} `json:"image,omitempty"`
	Source string      `json:"source,omitempty"`
}

// Comment represents a comment on a facebook post.
type Comment struct {
	ID                       string          `json:"id,omitempty"`
	Attachment               json.RawMessage `json:"attachment,omitempty"`
	CanComment               bool            `json:"can_comment,omitempty"`
	CanRemove                bool            `json:"can_remove,omitempty"`
	CanHide                  bool            `json:"can_hide,omitempty"`
	CanLike                  bool            `json:"can_like,omitempty"`
	CanReplyPrivately        bool            `json:"can_reply_privately,omitempty"`
	CommentCount             int32           `json:"comment_count,omitempty"`
	CreatedTime              fb.Time         `json:"created_time,omitempty"`
	From                     *User           `json:"from,omitempty"`
	LikeCount                int32           `json:"like_count,omitempty"`
	Message                  string          `json:"message,omitempty"`
	MessageTags              []MessageTag    `json:"message_tags,omitempty"`
	Object                   json.RawMessage `json:"object,omitempty"`
	Parent                   *Comment        `json:"parent,omitempty"`
	PrivateReplyConversation json.RawMessage `json:"private_reply_conversation,omitempty"`
	UserLikes                bool            `json:"user_likes,omitempty"`
}

// MessageTag represents a tagged user or site in a comment.
type MessageTag struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Length int    `json:"length,omitempty"`
}

// User represents a facebook user.
type User struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Reactions contains a reation and how often it was performed on an object.
type Reactions map[string]uint64
