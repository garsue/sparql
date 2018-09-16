package main

import (
	"context"
	"log"
	"time"

	"github.com/garsue/sparql"
)

func main() {
	cli, err := sparql.New("http://ja.dbpedia.org/sparql",
		sparql.MaxIdleConns(100),
		sparql.IdleConnTimeout(90*time.Second),
		sparql.Timeout(30*time.Second),
		sparql.Prefix("dbpj", "http://ja.dbpedia.org/resource/"),
		sparql.Prefix("dbp-owl", "http://dbpedia.org/ontology/"),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := cli.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	// Simple use case
	result, err := cli.Query(ctx, `select distinct * where { <http://ja.dbpedia.org/resource/東京都> ?p ?o . } LIMIT 10`)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", result)

	// Parameter
	result, err = cli.Query(ctx, `select * where { ?s <http://dbpedia.org/ontology/wikiPageID> $1 . } LIMIT 1`, sparql.Param{
		Ordinal: 1,
		Value:   1529557,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", result)

	// IRI parameter
	result, err = cli.Query(ctx, `select * where { $1 <http://ja.dbpedia.org/property/name> ?name . } LIMIT 10`, sparql.Param{
		Ordinal: 1,
		Value:   sparql.IRI("http://ja.dbpedia.org/resource/ももいろクローバーZ"),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", result)

	// With language tags
	result, err = cli.Query(ctx, `select * where { ?s <http://ja.dbpedia.org/property/name> $1 . } LIMIT 10`, sparql.Param{
		Ordinal: 1,
		Value: sparql.Literal{
			Value:       "ももいろクローバー",
			LanguageTag: "ja",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", result)

	// Typed literal with IRI
	result, err = cli.Query(ctx, `select * where { ?s <http://dbpedia.org/ontology/wikiPageLength> $1 . } LIMIT 1`, sparql.Param{
		Ordinal: 1,
		Value: sparql.Literal{
			Value:    76516,
			DataType: sparql.IRI("http://www.w3.org/2001/XMLSchema#nonNegativeInteger"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", result)

	// Typed literal with PrefixedName
	result, err = cli.Query(ctx, `select * where { ?s <http://dbpedia.org/ontology/birthYear> $1 . } LIMIT 1`, sparql.Param{
		Ordinal: 1,
		Value: sparql.Literal{
			Value:    "1995",
			DataType: sparql.PrefixedName("xsd:gYear"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", result)

	// Parameterized subject and object
	result, err = cli.Query(ctx, `select * where { $1 $2 ?o . } LIMIT 1000`,
		sparql.Param{
			Ordinal: 1,
			Value:   sparql.PrefixedName("dbpj:有安杏果"),
		},
		sparql.Param{
			Ordinal: 2,
			Value:   sparql.PrefixedName("dbp-owl:birthYear"),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", result)

	// Ask
	result, err = cli.Query(ctx, `ask { dbpj:有安杏果 dbp-owl:birthYear "1995"^^xsd:gYear . }`)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", result)
}
