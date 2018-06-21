package sql

import (
	"context"
	"database/sql/driver"
	"time"

	"github.com/garsue/go-sparql"
)

// Connector generates `driver.Conn` with a context.
type Connector struct {
	driver driver.Driver
	Name   string
}

// Connect returns `driver.Conn` with a context.
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	client, err := sparql.New(
		c.Name,
		sparql.MaxIdleConns(100),
		sparql.IdleConnTimeout(90*time.Second),
		sparql.Timeout(30*time.Second),
	)
	if err != nil {
		return nil, err
	}

	return &Conn{
		Client: client,
	}, nil
}

// Driver returns underlying `driver.Driver`.
func (c *Connector) Driver() driver.Driver {
	return c.driver
}
