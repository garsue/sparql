package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/garsue/sparql"
	_ "github.com/garsue/sparql/sql"
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

	rows, err := db.QueryContext(ctx, `select distinct * where { <http://ja.dbpedia.org/resource/東京都> ?p ?o .  } LIMIT 1`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var p, o sparql.IRI
		if err := rows.Scan(&p, &o); err != nil {
			log.Fatal(err)
		}
		log.Printf("%T %v %T %v", p, p, o, o)
	}

	if err := rows.Close(); err != nil {
		log.Fatal(err)
	}

	var o int
	if err := db.QueryRowContext(
		ctx,
		`select distinct * where { <http://ja.dbpedia.org/resource/東京都> <http://dbpedia.org/ontology/wikiPageID> ?o .  } LIMIT 1`,
	).Scan(&o); err != nil {
		log.Fatal(err)
	}
	log.Printf("%T %v", o, o)

	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}
