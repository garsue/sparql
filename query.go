package sparql

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Results is a part of a SPARQL query result json.
type Results struct {
	Bindings []Bindings `json:"bindings"`
}

// Bindings is a part of a SPARQL query result json.
type Bindings map[string]struct {
	Type     Type        `json:"type"`
	DataType URI         `json:"datatype"`
	XMLLang  string      `json:"xml:lang"`
	Value    interface{} `json:"value"`
}

// Type is the binding value type.
type Type string

// Types https://www.w3.org/TR/rdf-sparql-json-res/#variable-binding-results
const (
	TypeURI          Type = "uri"
	TypeLiteral      Type = "literal"
	TypeTypedLiteral Type = "typed-literal"
	TypeBlankNode    Type = "bnode"
)

// Query queries to the endpoint.
func (c *Client) Query(
	ctx context.Context,
	query string,
	params ...Param,
) (QueryResult, error) {
	request, err := c.request(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(request)
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
		return nil, fmt.Errorf(
			"SPARQL query error. status code: %d msg: %s",
			resp.StatusCode,
			errMsg,
		)
	}

	return DecodeXMLQueryResult(resp.Body)
}

func (c *Client) request(
	ctx context.Context,
	query string,
	params ...Param,
) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, c.Endpoint, nil)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)

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
		if _, err := b.WriteString(uri.Ref()); err != nil {
			return nil, err
		}
		if _, err := b.WriteString("\n"); err != nil {
			return nil, err
		}
	}

	// Replace parameters
	replacePairs := make([]string, 0, 2*len(params))
	for _, arg := range params {
		replacePairs = append(
			replacePairs,
			fmt.Sprintf("$%d", arg.Ordinal),
			arg.Serialize(),
		)
	}
	if _, err2 := b.WriteString(strings.
		NewReplacer(replacePairs...).
		Replace(query)); err2 != nil {
		return nil, err2
	}

	// Build the query
	built := b.String()
	c.Logger.Debug.Println(built)
	url := request.URL.Query()
	url.Set("query", built)
	url.Set("format", "application/sparql-results+xml")
	request.URL.RawQuery = url.Encode()
	return request, nil
}
