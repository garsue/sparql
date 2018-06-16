package sparql

import (
	"context"
	"database/sql/driver"
	"log"
	"net/http"
)

// Conn connects to a SPARQL source
type Conn struct {
	HttpClient http.Client
	Logger     *log.Logger
	Source     string
}

// Ping sends a HTTP HEAD request to the source
func (c *Conn) Ping(ctx context.Context) error {
	resp, err := c.HttpClient.Head(c.Source)
	c.Logger.Printf("ping %+v", resp)
	return err
}

// Prepare returns a prepared statement.
// TODO not implemented yet
func (*Conn) Prepare(query string) (driver.Stmt, error) {
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
