package sparql

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	Results Results `json:"results"`
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

func New(endpoint string) *Client {
	return &Client{
		HttpClient: http.Client{},
		Logger:     *logger.New(),
		Endpoint:   endpoint,
	}
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
	c.Logger.Debug.Printf("ping %+v", resp)

	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}
	return nil
}

// Query queries to the endpoint.
func (c *Client) Query(ctx context.Context, query string, args []driver.NamedValue) (*QueryResult, error) {
	request, err := http.NewRequest(http.MethodGet, c.Endpoint, nil)
	if err != nil {
		return nil, err
	}

	// TODO replace with parser combinator
	replacePairs := make([]string, 0, 2*len(args))
	for _, arg := range args {
		replacePairs = append(replacePairs, fmt.Sprintf("$%d", arg.Ordinal), fmt.Sprintf("%v", arg.Value))
	}

	url := request.URL.Query()
	url.Set("query", strings.NewReplacer(replacePairs...).Replace(query))
	url.Set("format", "json")
	request.URL.RawQuery = url.Encode()

	resp, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer c.Logger.LogCloseError(resp.Body)
	c.Logger.Debug.Printf("query context %+v\n", resp)

	var result QueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	c.Logger.Debug.Printf("query result %+v\n", result)

	return &result, nil
}
