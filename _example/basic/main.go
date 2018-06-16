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
}
