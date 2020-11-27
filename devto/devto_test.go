package devto

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testApiKey = "some-test-api-key"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the API client being tested
	client *API

	// server is a test HTTP server used to provide mock API responses
	server *httptest.Server
)

func setup(opts ...Option) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client, _ = New(testApiKey, opts...)
	client.BaseURL = server.URL
}

func teardown() {
	server.Close()
}

func TestHeaders(t *testing.T) {
	setup()

	mux.HandleFunc("/users/me", func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, "GET", request.Method)
		assert.Equal(t, testApiKey, request.Header.Get("api-key"))
		assert.Equal(t, "application/json", request.Header.Get("Content-Type"))
	})

	client.GetCurrentUser()

	teardown()

	// Test with extra custom headers on the client
	// TODO: test the priority of headers, the passed in ones should take priority over any defaults
	header := http.Header{}
	header.Add("X-Test-Header", "Some Header")

	setup(WithHeader(header))

	mux.HandleFunc("/users/me", func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, "Some Header", request.Header.Get("X-Test-Header"))
		assert.Equal(t, "GET", request.Method)
		assert.Equal(t, testApiKey, request.Header.Get("api-key"))
		assert.Equal(t, "application/json", request.Header.Get("Content-Type"))
	})

	client.GetCurrentUser()

	teardown()
}

func TestApiErrors(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/articles", func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writer.Write([]byte(`{
"error": "this is the error message",
"status": 422
}`))
	})

	_, err := client.CreateArticle(CreateArticle{})

	assert.EqualError(t, err, "error from CreateArticle: HTTP Status 422: this is the error message")
}
