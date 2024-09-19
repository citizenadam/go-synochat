package synochat

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

type ClientOpts struct {
	BaseURL    string
	HTTPClient *http.Client
}

type Client struct {
	BaseURL    *url.URL // Changed back to exported field
	httpClient *http.Client
}

func NewClient(baseURL string) (*Client, error) {
	return NewCustomClient(baseURL, nil)
}

func NewCustomClient(baseURL string, httpClient *http.Client) (*Client, error) {
	if strings.TrimSpace(baseURL) == "" {
		return nil, errors.New("url cannot be empty or just whitespaces")
	}

	u, err := url.ParseRequestURI(baseURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, errors.New("invalid URL provided")
	}

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &Client{
		httpClient: httpClient,
		BaseURL:    u,
	}, nil
}
