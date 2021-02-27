package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

type Client struct {
	EndpointURL *url.URL
	HTTPClient  *http.Client
	UserAgent   string
}

func NewClient(endpointURL string, httpClient *http.Client, userAgent string) (*Client, error) {
	parsedURL, err := url.ParseRequestURI(endpointURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse url: %s", endpointURL)
	}

	client := &Client{
		EndpointURL: parsedURL,
		HTTPClient:  httpClient,
		UserAgent:   userAgent,
	}
	return client, nil
}

func (client *Client) NewRequest(ctx context.Context, method string, subPath string, query string, body io.Reader) (*http.Request, error) {
	endpointURL := *client.EndpointURL
	endpointURL.Path = path.Join(client.EndpointURL.Path, subPath)
	endpointURL.RawQuery = query

	req, err := http.NewRequest(method, endpointURL.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", client.UserAgent)

	return req, nil
}

func (client *Client) DecodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}