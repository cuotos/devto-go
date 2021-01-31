package devto

import "net/http"

type Option func(*Client) error

func WithAPIURL(url string) Option {
	return func(client *Client) error {
		client.BaseURL = url
		return nil
	}
}

func WithHeader(header http.Header) Option {
	return func(client *Client) error {
		client.headers = header
		return nil
	}
}

func (c *Client) parseOptions(opts ...Option) error {
	for _, option := range opts {
		err := option(c)
		if err != nil {
			return err
		}
	}

	return nil
}
