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

type API struct {
	BaseURL    string
	APIKey     string
	httpClient *http.Client
	headers    http.Header
}

func newClient(opts ...Option) (*API, error) {
	api := &API{
		BaseURL: apiURL,
	}

	if err := api.parseOptions(opts...); err != nil {
		return nil, fmt.Errorf("options parsing failed: %w", err)
	}

	if api.httpClient == nil {
		api.httpClient = http.DefaultClient
	}

	return api, nil
}

func New(key string, opts ...Option) (*API, error) {
	if key == "" {
		return nil, errors.New(errEmptyCredentials)
	}

	api, err := newClient(opts...)
	if err != nil {
		return nil, err
	}

	api.APIKey = key

	return api, nil
}

func (a *API) makeRequest(method, uri string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return []byte{}, fmt.Errorf("error making request: %w", err)
	}

	// Merge the default headers (auth token) with headers passed in via the
	//  WithHeader option
	allHeaders := http.Header{}
	copyHeader(allHeaders, a.headers)

	req.Header = allHeaders

	req.Header.Set("api-key", a.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	//TODO: what if the response is paginated

	//TODO: handle >= 500 status codes etc

	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		respErr := Error{}
		if err := json.Unmarshal(respBody, &respErr); err != nil {
			return nil, Error{errUnmarshalError, resp.StatusCode}
		}
		return nil, Error{respErr.ErrorMsg, respErr.Status}
	}

	return respBody, nil
}

// copyHeader copies all headers for `source` and sets them on `target`.
// based on https://godoc.org/github.com/golang/gddo/httputil/header#Copy
func copyHeader(target, source http.Header) {
	for k, vs := range source {
		target[k] = vs
	}
}
