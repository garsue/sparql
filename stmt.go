package sparql

import (
	"context"
	"database/sql/driver"
)

// Stmt implements `driver.Stmt` with `sparql.Statement`.
type Stmt struct {
	*Statement
}

// QueryContext queries to a SPARQL source.
func (s *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	result, err := s.Statement.Query(ctx, argsToParams(args)...)
	if err != nil {
		return nil, err
	}

	return &Rows{
		queryResult: result,
	}, nil
}

// Close closes the statement. Actually do nothing.
func (s *Stmt) Close() error {
	return nil
}

// NumInput is not supported. Always return -1.
func (s *Stmt) NumInput() int {
	return -1
}

// Exec is not supported. DO NOT USE.
func (*Stmt) Exec(args []driver.Value) (driver.Result, error) {
	panic("not supported")
}

// Query queries to the endpoint.
func (s *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	panic("deprecated")
}
