package sparql

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

// Driver accesses SPARQL sources.
type Driver struct{}

// nolint: gochecknoinits
func init() {
	sql.Register("sparql", &Driver{})
}

// Open returns `driver.Conn`.
func (d *Driver) Open(name string) (driver.Conn, error) {
	connector, err := d.OpenConnector(name)
	if err != nil {
		return nil, err
	}
	return connector.Connect(context.Background())
}

// OpenConnector returns `driver.Connector`.
func (d *Driver) OpenConnector(name string) (driver.Connector, error) {
	return NewConnector(d, name), nil
}
