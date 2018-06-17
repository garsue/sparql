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
	conn := Conn{
		Client: sparql.New(c.Name, 100, 90*time.Second, 30*time.Second),
	}
	// Ping to get a keep-alive connection
	if err := conn.Client.Ping(ctx); err != nil {
		return nil, err
	}
	return &conn, nil
}

// Driver returns underlying `driver.Driver`.
func (c *Connector) Driver() driver.Driver {
	return c.driver
}
