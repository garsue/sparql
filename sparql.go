package sparql

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"net/http"
	"os"
)

// Driver accesses SPARQL sources.
type Driver struct {
	Logger *log.Logger
}

func init() {
	sql.Register("sparql", &Driver{
		Logger: log.New(os.Stdout, "[SPARQL:DEBUG]", log.LstdFlags),
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
