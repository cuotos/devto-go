package devto

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllUsersArticles(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/articles/me/all", func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set("content-type", "application/json")
		_, _ = writer.Write([]byte(`[{
			"type_of": "article",
			"id": 65656,
			"title": "The Title Of My Article",
			"description": "An amazing article about some stuff",
			"cover_image": "http://cutekittens.com/1.jpg",
			"published": true
		}]`))
	})

	expectedArticles := []Article{{
		TypeOf:      "article",
		ID:          65656,
		Title:       "The Title Of My Article",
		Description: "An amazing article about some stuff",
		CoverImage:  "http://cutekittens.com/1.jpg",
		Published:   true,
	}}

	articles, _, err := client.GetUsersArticles()
	if assert.NoError(t, err) {
		assert.Equal(t, expectedArticles, articles)
	}
}

func TestCanGetASingleArticle(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/articles/me/all", func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.Header().Set("content-type", "application/json")
		_, _ = writer.Write([]byte(`[
			{
				"type_of": "article",
				"id": 65656,
				"title": "The Title Of My Article"
			},
			{
				"type_of": "article",
				"id": 12345,
				"title": "Another Article",
				"published": true
			}
		]`))
	})

	expectedArticle := Article{
		TypeOf:    "article",
		ID:        12345,
		Title:     "Another Article",
		Published: true,
	}

	article, found, _, err := client.GetUserArticleByID(12345)
	if assert.NoError(t, err) {
		if assert.True(t, found, "article not found, but it should have been") {
			assert.Equal(t, expectedArticle, article)
		}
	}
}

func TestCreateArticle(t *testing.T) {
	setup()
	defer teardown()

	called := false

	mux.HandleFunc("/articles", func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		assert.Equal(t, "application/json", request.Header.Get("Content-Type"))

		postedBody, err := ioutil.ReadAll(request.Body)
		defer request.Body.Close()
		if err != nil {
			log.Fatal("unable to parse submitted body, this is not a broken test")
		}

		assert.Equal(t, `{"article":{"title":"Test Title"}}`, string(postedBody))

		called = true

		// Just returning the ID as that will have been created by the server
		writer.Write([]byte(`{
			"id": 98765
		}`))
	})

	testArticle := CreateArticle{
		Title:     "Test Title",
		Published: false,
	}
	createdArticle, _, err := client.CreateArticle(testArticle)

	if assert.NoError(t, err) {
		assert.Equal(t, 98765, createdArticle.ID)
		assert.True(t, called, "create article endpoint was not called")
	}
}

func TestCreateArticleThatAlreadyExists(t *testing.T) {
	setup()
	defer teardown()

	called := false

	mux.HandleFunc("/articles", func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		assert.Equal(t, "application/json", request.Header.Get("Content-Type"))

		writer.WriteHeader(http.StatusUnprocessableEntity)
		writer.Write([]byte(`{"error":"Body markdown has already been taken","status":422}`))

		called = true
	})

	testArticle := CreateArticle{
		Title:     "Test Title",
		Published: false,
	}
	_, _, err := client.CreateArticle(testArticle)

	assert.True(t, called)

	if assert.Error(t, err) {
		assert.Equal(t, "error from CreateArticle: HTTP Status 422: Body markdown has already been taken", err.Error())
	}
}

func TestUpdateArticle(t *testing.T) {
	setup()
	defer teardown()

	called := false

	mux.HandleFunc("/articles/123", func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPut, request.Method)
		assert.Equal(t, "application/json", request.Header.Get("Content-Type"))

		postedBody, err := ioutil.ReadAll(request.Body)
		defer request.Body.Close()
		if err != nil {
			log.Fatal("unable to parse submitted body, this is not a broken test")
		}

		assert.Equal(t, `{"article":{"title":"Update Title"}}`, string(postedBody))

		_, _ = writer.Write([]byte(`{}`))

		called = true
	})

	testArticle := CreateArticle{
		Title:     "Update Title",
		Published: false,
	}

	_, _, err := client.UpdateArticle(123, testArticle)

	assert.True(t, called)

	assert.NoError(t, err)
}
