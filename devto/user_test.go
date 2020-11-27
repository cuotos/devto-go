package devto

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/me", func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, "GET", request.Method)
		writer.Header().Set("content-type", "application/json")
		_, _ = writer.Write([]byte(`{
			"type_of":"user",	
			"id":12345,
			"username":"test-username",
			"name":"Test P User",
			"summary":"User Summary",
			"twitter_username":"showMeTheTwitter",
			"github_username":"GithubUsername",
			"website_url":"https://my-awsome-site.com",
			"location":"Bournemouth, UK",
			"joined_at":"Oct 28, 1904",
			"profile_image":"profile_url.jpeg"
		}`))
	})

	expectedUser := User{
		TypeOf:          "user",
		ID:              12345,
		Username:        "test-username",
		Name:            "Test P User",
		Summary:         "User Summary",
		TwitterUsername: "showMeTheTwitter",
		GithubUsername:  "GithubUsername",
		WebsiteURL:      "https://my-awsome-site.com",
		Location:        "Bournemouth, UK",
		JoinedAt:        "Oct 28, 1904",
		ProfileImage:    "profile_url.jpeg",
	}

	user, err := client.GetCurrentUser()
	if assert.NoError(t, err) {
		assert.Equal(t, expectedUser, user)
	}

}
