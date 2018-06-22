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

	var s, n, v interface{}
	if err := db.QueryRowContext(ctx, `
	SELECT ?s, ?n, ?v FROM <http://hojin-info.go.jp/graph/hojin>
	WHERE {
	?n ic:名称/ic:表記 ?v .
	FILTER regex(?v, "マネー")
	} LIMIT 1`).Scan(&s, &n, &v); err != nil {
		log.Fatal(err)
	}
	log.Println(s)
	log.Println(n)
	log.Println(v)

	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}
