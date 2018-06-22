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

	rows, err := db.QueryContext(ctx, `
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
