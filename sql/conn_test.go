package sql_test

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/garsue/sparql"
	ssql "github.com/garsue/sparql/sql"
)

func ExampleConn_QueryContext() {
	db, err := sql.Open("sparql", "http://ja.dbpedia.org/sparql")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	if err2 := db.PingContext(ctx); err2 != nil {
		panic(err2)
	}

	rows, err := db.QueryContext(
		ctx,
		"select distinct * where "+
			"{ <http://ja.dbpedia.org/resource/東京都> ?p ?o .  } LIMIT 1",
	)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var p, o sparql.IRI
		if err := rows.Scan(&p, &o); err != nil {
			panic(err)
		}
		log.Printf("%T %v %T %v", p, p, o, o)
	}

	if err := rows.Close(); err != nil {
		panic(err)
	}

	var o int
	if err := db.QueryRowContext(
		ctx,
		"select distinct * where "+
			"{ <http://ja.dbpedia.org/resource/東京都> "+
			"<http://dbpedia.org/ontology/wikiPageID> ?o .  } LIMIT 1",
	).Scan(&o); err != nil {
		panic(err)
	}
	log.Printf("%T %v", o, o)

	if err := db.Close(); err != nil {
		panic(err)
	}
}

func ExampleConn_QueryContext_hojin_info() {
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
		panic(err)
	}

	rows, err := db.QueryContext(ctx, `
		SELECT ?v FROM <http://hojin-info.go.jp/graph/hojin>
		WHERE {
		?n ic:名称/ic:表記 ?v .
		FILTER regex(?v, "マネー")
		} LIMIT 100`)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			panic(err)
		}
		log.Println(v)
	}

	if err := rows.Close(); err != nil {
		panic(err)
	}

	if err := db.Close(); err != nil {
		panic(err)
	}
}
