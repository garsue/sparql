package sparql

import (
	"context"
	"database/sql/driver"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/garsue/go-sparql/logger"
)

// Conn connects to a SPARQL source.
type Conn struct {
	HttpClient http.Client
	Logger     *logger.Logger
	Source     string
}

// QueryContext queries to a SPARQL source.
func (c *Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	request, err := http.NewRequest(http.MethodGet, c.Source, nil)
	if err != nil {
		return nil, err
	}

	replacePairs := make([]string, 0, 2*len(args))
	for _, arg := range args {
		replacePairs = append(replacePairs, fmt.Sprintf("$%d", arg.Ordinal), fmt.Sprintf("%v", arg.Value))
	}

	url := request.URL.Query()
	url.Set("query", strings.NewReplacer(replacePairs...).Replace(query))
	request.URL.RawQuery = url.Encode()

	resp, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer c.Logger.LogCloseError(resp.Body)
	c.Logger.Debug.Printf("query context %+v\n", resp)

	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		return nil, err
	}

	return nil, nil
}

// Ping sends a HTTP HEAD request to the source.
func (c *Conn) Ping(ctx context.Context) error {
	request, err := http.NewRequest(http.MethodHead, c.Source, nil)
	if err != nil {
		return err
	}

	resp, err := c.HttpClient.Do(request.WithContext(ctx))
	if err != nil {
		return err
	}
	defer c.Logger.LogCloseError(resp.Body)
	c.Logger.Debug.Printf("ping %+v", resp)
	if 200 <= resp.StatusCode && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("status code %d", resp.StatusCode)
}

// Prepare returns a prepared statement.
func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	panic("implement me")
}

// Close closes this connection but nothing to do.
func (*Conn) Close() error {
	return nil
}

// Begin is not supported. SPARQL does not have transactions.
func (*Conn) Begin() (driver.Tx, error) {
	panic("not supported")
}
