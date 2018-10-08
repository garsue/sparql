[![Build Status](https://travis-ci.org/garsue/sparql.svg?branch=master)](https://travis-ci.org/garsue/sparql)
[![codecov](https://codecov.io/gh/garsue/sparql/branch/master/graph/badge.svg)](https://codecov.io/gh/garsue/sparql)

# Go SQL driver for SPARQL

**SUPER EXPERIMENTAL**

[GoDoc](https://godoc.org/github.com/garsue/sparql)

A [SPARQL](https://www.w3.org/TR/sparql11-protocol/)-Driver for Go.
Including [database/sql](https://golang.org/pkg/database/sql/) implementation.

## Usage

See [examples](https://github.com/garsue/go-sparql/tree/master/_example).

## FAQ

Q: Can I use `?` for placeholders?
A: No. Please use `$1`, `$2`, `$3`, ... as placeholders.

Q: Which version's golang is supported?
A: go 1.11.x or later.