package devto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Article struct {
	TypeOf       string `json:"type_of"`
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	CoverImage   string `json:"cover_image"`
	Published    bool   `json:"published"`
	BodyMarkdown string `json:"body_markdown"`
	//PublishedAt            time.Time `json:"published_at"`
	//TagList                []string  `json:"tag_list"`
	//Slug                   string    `json:"slug"`
	//Path                   string    `json:"path"`
	//URL                    string    `json:"url"`
	//CanonicalURL           string    `json:"canonical_url"`
	//CommentsCount          int       `json:"comments_count"`
	//PositiveReactionsCount int       `json:"positive_reactions_count"`
	//PublicReactionsCount   int       `json:"public_reactions_count"`
	//PageViewsCount         int       `json:"page_views_count"`
	//PublishedTimestamp     time.Time `json:"published_timestamp"`
	//User                   User      `json:"user"`
	//Organization           struct {
	//	Name           string `json:"name"`
	//	Username       string `json:"username"`
	//	Slug           string `json:"slug"`
	//	ProfileImage   string `json:"profile_image"`
	//	ProfileImage90 string `json:"profile_image_90"`
	//} `json:"organization"`
	//FlareTag struct {
	//	Name         string `json:"name"`
	//	BgColorHex   string `json:"bg_color_hex"`
	//	TextColorHex string `json:"text_color_hex"`
	//} `json:"flare_tag"`
}

type CreateArticle struct {
	Title        string   `json:"title"`
	Published    bool     `json:"published,omitempty"`
	BodyMarkdown string   `json:"body_markdown,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Series       string   `json:"series,omitempty"`
	CanonicalURL string   `json:"canonical_url,omitempty"`
}

// GetUsersArticles returns a slice of Articles for the authenticated user.
//
// API reference: https://docs.dev.to/api/#operation/getUserArticles
func (a *API) GetUsersArticles() ([]Article, error) {

	uri := "/articles/me/all"

	resp, err := a.makeRequest(http.MethodGet, a.BaseURL+uri, nil)
	if err != nil {
		return []Article{}, err
	}

	articles := []Article{}
	err = json.Unmarshal(resp, &articles)
	if err != nil {
		return []Article{}, err
	}

	return articles, nil
}

// GetUserArticleByID returns a single article
//
// As no api exists to get a single unpublished article, this gets all of the users articles and filters the required one
func (a *API) GetUserArticleByID(id int) (Article, bool, error) {
	articles, err := a.GetUsersArticles()
	if err != nil {
		return Article{}, false, err
	}

	for _, a := range articles {
		if a.ID == id {
			return a, true, nil
		}
	}

	return Article{}, false, nil
}

//CreateArticle creates an article for the currently authenticated user
//
//API reference: https://docs.dev.to/api/#operation/createArticle
func (a *API) CreateArticle(article CreateArticle) (Article, error) {
	return a.upsert(nil, article)
}

//UpdateArticle updates an existing article owned by the currently authenticated user
//
//API Reference: https://docs.dev.to/api/#operation/updateArticle
func (a *API) UpdateArticle(id int, article CreateArticle) (Article, error) {
	return a.upsert(&id, article)
}

func (a *API) upsert(id *int, article CreateArticle) (Article, error) {
	var uri string
	var method string

	if id != nil {
		uri = fmt.Sprintf("/articles/%d", *id)
		method = http.MethodPut
	} else {
		uri = "/articles"
		method = http.MethodPost
	}

	type CreateArticleRequest struct {
		Article CreateArticle `json:"article"`
	}

	createArticleRequest := CreateArticleRequest{
		Article: article,
	}

	articleBytes, err := json.Marshal(createArticleRequest)
	if err != nil {
		return Article{}, fmt.Errorf("unable to marshal article: %w", err)
	}

	resp, err := a.makeRequest(method, a.BaseURL+uri, bytes.NewBuffer(articleBytes))
	if err != nil {
		return Article{}, fmt.Errorf("error from CreateArticle: HTTP Status %v: %v", err.(Error).Status, err.Error())
	}

	createArticleResponse := Article{}

	err = json.Unmarshal(resp, &createArticleResponse)
	if err != nil {
		return Article{}, fmt.Errorf("unable to unmarshal the create article response object: %w", err)
	}

	return createArticleResponse, nil
}
