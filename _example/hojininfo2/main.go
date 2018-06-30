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

	// MF
	rows, err := db.QueryContext(ctx, `
		SELECT ?v FROM <http://hojin-info.go.jp/graph/hojin>
		WHERE {
		?n ic:名称/ic:表記 ?v .
		FILTER regex(?v, "マネー")
		} LIMIT 100`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			log.Fatal(err)
		}
		log.Println(v)
	}

	if err := rows.Close(); err != nil {
		log.Fatal(err)
	}

	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}
