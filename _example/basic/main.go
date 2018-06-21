package main

import (
	"context"
	"log"
	"time"

	"github.com/garsue/go-sparql"
)

func main() {
	cli := sparql.New("http://ja.dbpedia.org/sparql", 100, 90*time.Second, 30*time.Second)

	ctx := context.Background()
	if err := cli.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	result, err := cli.Query(ctx, `select distinct * where { <http://ja.dbpedia.org/resource/東京都> ?p ?o . } LIMIT 10`)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", result)

	result, err = cli.Query(ctx, `select * where { ?s <http://dbpedia.org/ontology/wikiPageID> $1 . } LIMIT 1`, sparql.Param{
		Ordinal: 1,
		Value:   1529557,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", result)

	result, err = cli.Query(ctx, `select * where { $1 <http://ja.dbpedia.org/property/name> ?name . } LIMIT 10`, sparql.Param{
		Ordinal: 1,
		Value:   sparql.SparqlURL("http://ja.dbpedia.org/resource/ももいろクローバーZ"),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", result)

	result, err = cli.Query(ctx, `select * where { ?s <http://ja.dbpedia.org/property/name> $1 . } LIMIT 10`, sparql.Param{
		Ordinal:     1,
		Value:       "ももいろクローバー",
		LanguageTag: "ja",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", result)
}
