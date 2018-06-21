package sparql

import (
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
