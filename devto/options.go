package devto

import "net/http"

type Option func(*API) error

func WithAPIURL(url string) Option {
	return func(api *API) error {
		api.BaseURL = url
		return nil
	}
}

func WithHeader(header http.Header) Option {
	return func(api *API) error {
		api.headers = header
		return nil
	}
}

func (a *API) parseOptions(opts ...Option) error {
	for _, option := range opts {
		err := option(a)
		if err != nil {
			return err
		}
	}

	return nil
}
