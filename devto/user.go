package devto

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	TypeOf          string      `json:"type_of"`
	ID              int         `json:"id"`
	Username        string      `json:"username"`
	Name            string      `json:"name"`
	Summary         string      `json:"summary"`
	TwitterUsername interface{} `json:"twitter_username"`
	GithubUsername  string      `json:"github_username"`
	WebsiteURL      string      `json:"website_url"`
	Location        string      `json:"location"`
	JoinedAt        string      `json:"joined_at"`
	ProfileImage    string      `json:"profile_image"`
}

// GetCurrentUser returns the currently authenticated User.
//
// API reference: https://docs.dev.to/api/#operation/getUserMe
func (a *API) GetCurrentUser() (User, error) {

	uri := "/users/me"

	resp, err := a.makeRequest(http.MethodGet, a.BaseURL+uri, nil)
	if err != nil {
		return User{}, err
	}

	user := User{}
	err = json.Unmarshal(resp, &user)
	if err != nil {
		return User{}, fmt.Errorf("unable to unmarshall response: %w", err)
	}

	return user, nil
}
