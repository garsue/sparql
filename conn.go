package sparql

import (
	"context"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/garsue/sparql/client"
)

// Conn connects to a SPARQL source.
type Conn struct {
	Client *client.Client
}

// Rows implements `driver.Rows` with `sparql.QueryResult`.
type Rows struct {
	queryResult client.QueryResult
}

// Columns returns the names of the columns.
func (r *Rows) Columns() []string {
	return r.queryResult.Variables()
}

// Close closes the rows iterator.
func (r *Rows) Close() error {
	return r.queryResult.Close()
}

// Next is called to populate the next row of data into
// the provided slice.
func (r *Rows) Next(dest []driver.Value) error {
	// See Boolean field if the query is ASK (not SELECT)
	variables := r.queryResult.Variables()
	if len(variables) == 0 && len(dest) > 0 {
		b, err := r.queryResult.Boolean()
		if err != nil {
			return err
		}
		dest[0] = b
		return nil
	}

	bindings, err := r.queryResult.Next()
	if err != nil {
		return err
	}
	for i, k := range variables {
		dest[i] = scan(bindings[k])
	}
	return nil
}

func scan(b client.Value) driver.Value {
	literal, ok := b.(client.Literal)
	if !ok {
		return b
	}
	if literal.DataType == client.URI("http://www.w3.org/2001/XMLSchema#dateTime") {
		for _, f := range []string{
			"2006-01-02T15:04:05.999999999",
			time.RFC3339Nano,
		} {
			t, err := time.ParseInLocation(f, fmt.Sprint(literal.Value), time.UTC)
			if err == nil {
				return t
			}
		}
	}
	return literal.Value
}

// QueryContext queries to a SPARQL source.
func (c *Conn) QueryContext(
	ctx context.Context,
	query string,
	args []driver.NamedValue,
) (driver.Rows, error) {
	result, err := c.Client.Query(ctx, query, argsToParams(args)...)
	if err != nil {
		return nil, err
	}

	return &Rows{
		queryResult: result,
	}, nil
}

func argsToParams(args []driver.NamedValue) []client.Param {
	params := make([]client.Param, 0, len(args))
	for _, a := range args {
		params = append(params, client.Param{
			Name:    a.Name,
			Ordinal: a.Ordinal,
			Value:   interface{}(a.Value),
		})
	}
	return params
}

// Ping sends a HTTP HEAD request to the source.
func (c *Conn) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx)
}

// Prepare returns a prepared statement.
func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	return &Stmt{
		Statement: c.Client.Prepare(query),
	}, nil
}

// Close closes this connection but nothing to do.
func (c *Conn) Close() error {
	return c.Client.Close()
}

// Begin is not supported. SPARQL does not have transactions.
func (*Conn) Begin() (driver.Tx, error) {
	panic("not supported")
}
