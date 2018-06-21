package sql

import (
	"context"
	"database/sql/driver"

	"github.com/garsue/go-sparql"
)

// Conn connects to a SPARQL source.
type Conn struct {
	Client *sparql.Client
}

// Query queries to a SPARQL source.
func (c *Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	params := make([]sparql.Param, 0, len(args))
	for _, a := range args {
		params = append(params, sparql.Param{
			Ordinal: a.Ordinal,
			Value:   a.Value.(interface{}),
		})
	}

	_, err := c.Client.Query(ctx, query, params...)
	return nil, err
}

// Ping sends a HTTP HEAD request to the source.
func (c *Conn) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx)
}

// Prepare returns a prepared statement.
func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	panic("implement me")
}

// Close closes this connection but nothing to do.
func (c *Conn) Close() error {
	return c.Client.Close()
}

// Begin is not supported. SPARQL does not have transactions.
func (*Conn) Begin() (driver.Tx, error) {
	panic("not supported")
}
