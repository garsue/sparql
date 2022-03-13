[![GoDoc](https://godoc.org/github.com/garsue/sparql?status.svg)](https://godoc.org/github.com/garsue/sparql)
![Build Status](https://github.com/garsue/sparql/actions/workflows/go.yml/badge.svg)
[![codecov](https://codecov.io/gh/garsue/sparql/branch/master/graph/badge.svg)](https://codecov.io/gh/garsue/sparql)
[![Go Report Card](https://goreportcard.com/badge/github.com/garsue/sparql)](https://goreportcard.com/report/github.com/garsue/sparql)

# Go SQL driver for SPARQL

**SUPER EXPERIMENTAL**

A [SPARQL](https://www.w3.org/TR/sparql11-protocol/)-Driver for Go.
Including [database/sql](https://golang.org/pkg/database/sql/) implementation.

## Usage

See [examples](https://github.com/garsue/go-sparql/tree/master/_example).

## FAQ

Q: Can I use `?` for placeholders?
A: No. Please use `$1`, `$2`, `$3`, ... as placeholders.

Q: Which version's golang is supported?
A: go 1.11.x or later.
