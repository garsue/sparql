package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/garsue/go-sparql/sql"
)

func main() {
	db, err := sql.Open("sparql", "https://api.hojin-info.go.jp/sparql")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	rows, err := db.QueryContext(ctx, `PREFIX hj: <http://hojin-info.go.jp/ns/domain/biz/1#>
PREFIX ic: <http://imi.go.jp/ns/core/rdf#>
SELECT DISTINCT ?cID (SAMPLE(?name) AS ?name) (SUM(?value) AS ?sum)
FROM <http://hojin-info.go.jp/graph/chotatsu>
WHERE {
?s hj:法人活動情報 ?o .
?o ic:ID/ic:識別値 ?cID .
?o ic:名称/ic:表記 ?name .
?o ic:金額/ic:数値 ?value .
}
GROUP BY ?cID
ORDER BY DESC(?sum)`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var cID, name string
		var sum int64
		if err := rows.Scan(&cID, &name, &sum); err != nil {
			log.Fatal(err)
		}
		log.Println(cID, name, sum)
	}

	if err := rows.Close(); err != nil {
		log.Fatal(err)
	}

	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}
