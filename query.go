package sparql

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Query queries to the endpoint.
func (c *Client) Query(
	ctx context.Context,
	query string,
	params ...Param,
) (QueryResult, error) {
	return c.Prepare(query).Query(ctx, params...)
}

// Statement is prepared statement.
type Statement struct {
	c      *Client
	query  string
	prefix string
}

// Prepare returns `*sparql.Statement`.
func (c *Client) Prepare(query string) *Statement {
	// Prepare PREFIX
	ss := make([]string, 0, len(c.prefixes)*5)
	for prefix, uri := range c.prefixes {
		ss = append(ss, "PREFIX ")
		ss = append(ss, prefix)
		ss = append(ss, ": ")
		ss = append(ss, uri.Ref())
		ss = append(ss, "\n")
	}
	return &Statement{c: c, prefix: strings.Join(ss, ""), query: query}
}

// Query queries to the endpoint.
func (s *Statement) Query(
	ctx context.Context,
	params ...Param,
) (QueryResult, error) {
	request, err := s.request(ctx, params...)
	if err != nil {
		return nil, err
	}

	resp, err := s.c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer s.c.Logger.LogCloseError(resp.Body)
	s.c.Logger.Debug.Printf("%+v\n", resp)

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

func (s *Statement) request(ctx context.Context, params ...Param) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, s.c.Endpoint, nil)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)

	const defaultBufferSize = 1024
	b := bytes.NewBuffer(make([]byte, 0, defaultBufferSize))

	if err := s.compose(b, params...); err != nil {
		return nil, err
	}

	// Build the query
	built := b.String()
	s.c.Logger.Debug.Println(built)
	url := request.URL.Query()
	url.Set("query", built)
	url.Set("format", "application/sparql-results+xml")
	request.URL.RawQuery = url.Encode()
	return request, nil
}

func (s *Statement) compose(writer io.Writer, params ...Param) error {
	// Write prefix
	if _, err := writer.Write([]byte(s.prefix)); err != nil {
		return err
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

	_, err := strings.NewReplacer(replacePairs...).WriteString(writer, s.query)
	return err
}
