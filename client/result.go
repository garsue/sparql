package client

import (
	"io"
)

// ResultParser is the parser for specific format query results.
type ResultParser interface {
	// Format returns a format name string. It's used as a `format` request header value.
	Format() string
	// Parse parses query result stream.
	Parse(reader io.ReadCloser) (QueryResult, error)
}

// QueryResult is a SPARQL query result.
type QueryResult interface {
	Variables() []string
	Next() (map[string]Value, error)
	Boolean() (bool, error)

	io.Closer
}

// Value is an interface holding one of the binding (or boolean) types:
// URI, Literal, BNode or bool.
type Value interface{}
