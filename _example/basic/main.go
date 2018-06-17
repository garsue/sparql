package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/garsue/go-sparql"
)

func main() {
	db, err := sql.Open("sparql", "http://ja.dbpedia.org/sparql")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	rows, err := db.QueryContext(ctx, `select distinct * where { <http://ja.dbpedia.org/resource/東京都> ?p ?o .  } LIMIT 100`)
	if err != nil {
		log.Fatal(err)
	}
	if err := rows.Close(); err != nil {
		log.Fatal(err)
	}

	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}
