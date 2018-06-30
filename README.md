# Go SQL driver for SPARQL

**SUPER EXPERIMENTAL**

A [SPARQL](https://www.w3.org/TR/sparql11-protocol/)-Driver for Go's [database/sql](https://golang.org/pkg/database/sql/) package.

## Usage

See [_example](https://github.com/garsue/go-sparql/tree/master/_example)

## FAQ

Q: Can I use `?` for placeholders?
A: No. Please use `$1`, `$2`, `$3`, ... as placeholders.

Q: Which version's golang is supported?
A: go 1.10.x or laterj