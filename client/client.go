package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/garsue/sparql/logger"
)

// Client queries to its SPARQL endpoint.
type Client struct {
	HTTPClient http.Client
	Logger     logger.Logger
	Endpoint   string
	prefixes   map[string]URI
}

// Option sets an option to the SPARQL client.
type Option func(*Client) error

// HTTPClient replaces default HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) error {
		c.HTTPClient = *httpClient
		return nil
	}
}

// WithPrefix sets a global PREFIX for all queries.
func WithPrefix(prefix string, uri URI) Option {
	return func(c *Client) error {
		c.prefixes[prefix] = uri
		return nil
	}
}

// New returns `sparql.Client`.
func New(endpoint string, opts ...Option) (*Client, error) {
	client := &Client{
		Logger:   *logger.New(),
		Endpoint: endpoint,
		prefixes: make(map[string]URI),
	}
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}
	return client, nil
}

// Close closes this client.
// Actually nothing to do to close the HTTP client.
func (c *Client) Close() error {
	return nil
}

// Ping sends a HTTP HEAD request to the endpoint.
func (c *Client) Ping(ctx context.Context) error {
	request, err := http.NewRequest(http.MethodHead, c.Endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(request.WithContext(ctx))
	if err != nil {
		return err
	}
	defer c.Logger.LogCloseError(resp.Body)
	c.Logger.Debug.Printf("Ping %+v", resp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SPARQL ping error. status code %d", resp.StatusCode)
	}
	return nil
}
