package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/garsue/go-sparql"
	ssql "github.com/garsue/go-sparql/sql"
)

func main() {
	db := sql.OpenDB(ssql.NewConnector(
		"https://api.hojin-info.go.jp/sparql",
		sparql.Timeout(5*time.Second),
		sparql.MaxIdleConns(0),
		sparql.IdleConnTimeout(0),
		sparql.Prefix("hj", "http://hojin-info.go.jp/ns/domain/biz/1#"),
		sparql.Prefix("ic", "http://imi.go.jp/ns/core/rdf#"),
	))

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	var b bool
	if err := db.QueryRowContext(ctx, `
	ASK FROM <http://hojin-info.go.jp/graph/hyosho>
	WHERE {
	?n ic:名称/ic:表記 ?v .
	FILTER regex(?v, "マネー")
	}`).Scan(&b); err != nil {
		log.Fatal(err)
	}
	log.Println(b)

	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}
