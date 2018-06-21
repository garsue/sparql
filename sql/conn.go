package sql

import (
	"context"
	"database/sql/driver"
	"io"

	"github.com/garsue/go-sparql"
)

// Conn connects to a SPARQL source.
type Conn struct {
	Client *sparql.Client
}

// Rows implements `driver.Rows` with `sparql.QueryResult`.
type Rows struct {
	queryResult *sparql.QueryResult
}

// Columns retunrs the names of the columns.
func (r *Rows) Columns() []string {
	return r.queryResult.Head.Vars
}

// Close closes the rows iterator.
func (r *Rows) Close() error {
	r.queryResult = nil
	return nil
}

// Next is called to populate the next row of data into
// the provided slice.
func (r *Rows) Next(dest []driver.Value) error {
	for _, b := range r.queryResult.Results.Bindings {
		for i, k := range r.queryResult.Head.Vars {
			dest[i] = b[k]
		}
		r.queryResult.Results.Bindings = r.queryResult.Results.Bindings[1:]
		return nil
	}
	return io.EOF
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

	result, err := c.Client.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	return &Rows{
		queryResult: result,
	}, err
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
