package devto

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	apiURL = "https://dev.to/api/"
)

type Client struct {
	BaseURL    string
	APIKey     string
	httpClient *http.Client
	headers    http.Header
}

func newClient(opts ...Option) (*Client, error) {
	client := &Client{
		BaseURL: apiURL,
	}

	if err := client.parseOptions(opts...); err != nil {
		return nil, fmt.Errorf("options parsing failed: %w", err)
	}

	if client.httpClient == nil {
		client.httpClient = http.DefaultClient
	}

	return client, nil
}

func New(key string, opts ...Option) (*Client, error) {
	if key == "" {
		return nil, errors.New(errEmptyCredentials)
	}

	client, err := newClient(opts...)
	if err != nil {
		return nil, err
	}

	client.APIKey = key

	return client, nil
}

func (c *Client) makeRequest(method, uri string, body io.Reader) ([]byte, *Response, error) {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return []byte{}, nil, fmt.Errorf("error making request: %w", err)
	}

	// Merge the default headers (auth token) with headers passed in via the
	//  WithHeader option
	allHeaders := http.Header{}
	copyHeader(allHeaders, c.headers)

	req.Header = allHeaders

	req.Header.Set("api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(req)
	resp := &Response{response}
	if err != nil {
		return nil, resp, fmt.Errorf("HTTP request failed: %w", err)
	}

	//TODO: what if the response is paginated

	//TODO: handle >= 500 status codes etc

	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, resp, fmt.Errorf("could not read response body: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		respErr := Error{}
		if err := json.Unmarshal(respBody, &respErr); err != nil {
			return nil, resp, errors.New(errUnmarshalError)
		}
		return nil, resp, err
	}

	return respBody, resp, nil
}

// copyHeader copies all headers for `source` and sets them on `target`.
// based on https://godoc.org/github.com/golang/gddo/httputil/header#Copy
func copyHeader(target, source http.Header) {
	for k, vs := range source {
		target[k] = vs
	}
}

// Response is a DEV.to API reponse. This wraps the standard http.Response
// allowing the addition of helper functions in the future.
type Response struct {
	*http.Response
}
