package sparql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/garsue/go-sparql/logger"
)

// Client queries to its SPARQL endpoint.
type Client struct {
	HttpClient http.Client
	Logger     logger.Logger
	Endpoint   string
}

// QueryResult is a destination to decoding a SPARQL query result json.
type QueryResult struct {
	Head    Head    `json:"head"`
	Results Results `json:"results"`
}

// Head is a part of a SPARQL query result json.
type Head struct {
	Vars []string `json:"vars"`
}

// Results is a part of a SPARQL query result json.
type Results struct {
	Bindings []Binding `json:"bindings"`
}

// Binding is a part of a SPARQL query result json.
type Binding map[string]struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// New returns `sparql.Client`.
func New(endpoint string, maxIdleConns int, idleConnTimeout, timeout time.Duration) *Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: timeout,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          maxIdleConns,
		MaxIdleConnsPerHost:   maxIdleConns,
		IdleConnTimeout:       idleConnTimeout,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &Client{
		HttpClient: http.Client{
			Transport: transport,
		},
		Logger:   *logger.New(),
		Endpoint: endpoint,
	}
}

// Close closes this client
func (c *Client) Close() error {
	if c.HttpClient.Transport == nil {
		return errors.New("already closed")
	}
	transport, ok := c.HttpClient.Transport.(*http.Transport)
	if !ok {
		return fmt.Errorf("unknown RoundTripper, %+v", c.HttpClient.Transport)
	}
	transport.CloseIdleConnections()
	c.HttpClient.Transport = nil
	return nil
}

// Ping sends a HTTP HEAD request to the endpoint.
func (c *Client) Ping(ctx context.Context) error {
	request, err := http.NewRequest(http.MethodHead, c.Endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := c.HttpClient.Do(request.WithContext(ctx))
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

// Query queries to the endpoint.
func (c *Client) Query(ctx context.Context, query string, args ...Param) (*QueryResult, error) {
	request, err := http.NewRequest(http.MethodGet, c.Endpoint, nil)
	if err != nil {
		return nil, err
	}

	replacePairs := make([]string, 0, 2*len(args))
	for _, arg := range args {
		replacePairs = append(replacePairs, fmt.Sprintf("$%d", arg.Ordinal), arg.Serialize())
	}

	url := request.URL.Query()
	built := strings.NewReplacer(replacePairs...).Replace(query)
	c.Logger.Debug.Println(built)
	url.Set("query", built)
	url.Set("format", "application/sparql-results+json")
	request.URL.RawQuery = url.Encode()

	resp, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer c.Logger.LogCloseError(resp.Body)
	c.Logger.Debug.Printf("%+v\n", resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SPARQL query error. status code %d", resp.StatusCode)
	}

	var result QueryResult
	//body := io.TeeReader(resp.Body, os.Stdout)
	body := resp.Body
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
