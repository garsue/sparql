package sparql

import (
	"context"
	"database/sql/driver"

	"github.com/garsue/sparql/client"
)

// Connector generates `driver.Conn` with a context.
type Connector struct {
	driver  driver.Driver
	Name    string
	options []client.Option
}

// NewConnector returns `driver.Connector`.
func NewConnector(
	driver driver.Driver,
	name string,
	opts ...client.Option,
) *Connector {
	return &Connector{
		driver:  driver,
		Name:    name,
		options: opts,
	}
}

// Connect returns `driver.Conn` with a context.
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	cli, err := client.New(c.Name, c.options...)
	if err != nil {
		return nil, err
	}

	return &Conn{
		Client: cli,
	}, nil
}

// Driver returns underlying `driver.Driver`.
func (c *Connector) Driver() driver.Driver {
	return c.driver
}
