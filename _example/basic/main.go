package main

import (
	"context"
	"log"

	"github.com/garsue/go-sparql"
)

func main() {
	cli := sparql.New("http://ja.dbpedia.org/sparql")

	ctx := context.Background()
	if err := cli.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	result, err := cli.Query(ctx, `select distinct * where { <http://ja.dbpedia.org/resource/東京都> ?p ?o .  } LIMIT 100`, nil)
	if err != nil {
		log.Fatal(err)
	}
	if result == nil {
		log.Fatal("result should not be nil")
	}
}
