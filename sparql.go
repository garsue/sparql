package sparql

import (
	"database/sql"
	"database/sql/driver"
	"net/http"

	"github.com/garsue/go-sparql/logger"
)

// Driver accesses SPARQL sources.
type Driver struct {
	Logger *logger.Logger
}

func init() {
	sql.Register("sparql", &Driver{
		Logger: logger.New(),
	})
}

// Open opens a SPARQL source.
func (d *Driver) Open(name string) (driver.Conn, error) {
	return &Conn{
		HttpClient: http.Client{},
		Logger:     d.Logger,
		Source:     name,
	}, nil
}
