package sparql

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// QueryResult is a destination to decoding a SPARQL query result json.
type QueryResult struct {
	Head    Head    `json:"head"`
	Results Results `json:"results"`
	Boolean bool    `json:"boolean"`
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

// Query queries to the endpoint.
func (c *Client) Query(ctx context.Context, query string, args ...Param) (*QueryResult, error) {
	request, err := http.NewRequest(http.MethodGet, c.Endpoint, nil)
	if err != nil {
		return nil, err
	}

	const defaultBufferSize = 1024
	b := bytes.NewBuffer(make([]byte, 0, defaultBufferSize))

	// Prepend PREFIX
	for prefix, uri := range c.prefixes {
		if _, err := b.WriteString("PREFIX "); err != nil {
			return nil, err
		}
		if _, err := b.WriteString(prefix); err != nil {
			return nil, err
		}
		if _, err := b.WriteString(": "); err != nil {
			return nil, err
		}
		if _, err := b.WriteString(uri.Serialize()); err != nil {
			return nil, err
		}
		if _, err := b.WriteString("\n"); err != nil {
			return nil, err
		}
	}

	// Replace parameters
	replacePairs := make([]string, 0, 2*len(args))
	for _, arg := range args {
		replacePairs = append(replacePairs, fmt.Sprintf("$%d", arg.Ordinal), arg.Serialize())
	}
	if _, err := b.WriteString(strings.NewReplacer(replacePairs...).Replace(query)); err != nil {
		return nil, err
	}

	// Build the query
	built := b.String()
	c.Logger.Debug.Println(built)
	url := request.URL.Query()
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
		scanner := bufio.NewScanner(resp.Body)
		var errMsg string
		if scanner.Scan() {
			errMsg = scanner.Text()
		}
		return nil, fmt.Errorf("SPARQL query error. status code: %d msg: %s", resp.StatusCode, errMsg)
	}

	var result QueryResult
	//body := io.TeeReader(resp.Body, os.Stdout)
	body := resp.Body
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
