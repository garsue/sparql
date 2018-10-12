package sparql

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/garsue/sparql/logger"
)

func ExampleClient_Query_simple() {
	cli, err := New("http://ja.dbpedia.org/sparql",
		MaxIdleConns(100),
		IdleConnTimeout(90*time.Second),
		Timeout(30*time.Second),
		Prefix("dbpj", "http://ja.dbpedia.org/resource/"),
		Prefix("dbp-owl", "http://dbpedia.org/ontology/"),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	if err2 := cli.Ping(ctx); err2 != nil {
		panic(err2)
	}

	// Simple use case
	result, err := cli.Query(
		ctx,
		"select distinct * where "+
			"{ <http://ja.dbpedia.org/resource/東京都> ?p ?o . } LIMIT 10",
	)
	if err != nil {
		panic(err)
	}
	log.Println(result.Next())
	if err := result.Close(); err != nil {
		panic(err)
	}
}

func ExampleClient_Query_parameter() {
	cli, err := New("http://ja.dbpedia.org/sparql",
		MaxIdleConns(100),
		IdleConnTimeout(90*time.Second),
		Timeout(30*time.Second),
		Prefix("dbpj", "http://ja.dbpedia.org/resource/"),
		Prefix("dbp-owl", "http://dbpedia.org/ontology/"),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	if err2 := cli.Ping(ctx); err2 != nil {
		panic(err2)
	}

	// Parameter
	result, err := cli.Query(
		ctx,
		`select * where { ?s <http://dbpedia.org/ontology/wikiPageID> $1 . } LIMIT 1`,
		Param{
			Ordinal: 1,
			Value:   1529557,
		},
	)
	if err != nil {
		panic(err)
	}
	log.Println(result.Next())
	if err := result.Close(); err != nil {
		panic(err)
	}
}

func ExampleClient_Query_uri_parameter() {
	cli, err := New("http://ja.dbpedia.org/sparql",
		MaxIdleConns(100),
		IdleConnTimeout(90*time.Second),
		Timeout(30*time.Second),
		Prefix("dbpj", "http://ja.dbpedia.org/resource/"),
		Prefix("dbp-owl", "http://dbpedia.org/ontology/"),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	if err2 := cli.Ping(ctx); err2 != nil {
		panic(err2)
	}

	// URI parameter
	result, err := cli.Query(
		ctx,
		"select * where "+
			"{ $1 <http://ja.dbpedia.org/property/name> ?name . } LIMIT 10",
		Param{
			Ordinal: 1,
			Value:   URI("http://ja.dbpedia.org/resource/ももいろクローバーZ"),
		},
	)
	if err != nil {
		panic(err)
	}
	log.Println(result.Next())
	if err := result.Close(); err != nil {
		panic(err)
	}
}

func ExampleClient_Query_language_tag() {
	cli, err := New("http://ja.dbpedia.org/sparql",
		MaxIdleConns(100),
		IdleConnTimeout(90*time.Second),
		Timeout(30*time.Second),
		Prefix("dbpj", "http://ja.dbpedia.org/resource/"),
		Prefix("dbp-owl", "http://dbpedia.org/ontology/"),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	if err2 := cli.Ping(ctx); err2 != nil {
		panic(err2)
	}

	// With language tags
	result, err := cli.Query(
		ctx,
		`select * where { ?s <http://ja.dbpedia.org/property/name> $1 . } LIMIT 10`,
		Param{
			Ordinal: 1,
			Value: Literal{
				Value:       "ももいろクローバー",
				LanguageTag: "ja",
			},
		},
	)
	if err != nil {
		panic(err)
	}
	log.Println(result.Next())
	if err := result.Close(); err != nil {
		panic(err)
	}
}

func ExampleClient_Query_typed_literal_with_uri() {
	cli, err := New("http://ja.dbpedia.org/sparql",
		MaxIdleConns(100),
		IdleConnTimeout(90*time.Second),
		Timeout(30*time.Second),
		Prefix("dbpj", "http://ja.dbpedia.org/resource/"),
		Prefix("dbp-owl", "http://dbpedia.org/ontology/"),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	if err2 := cli.Ping(ctx); err2 != nil {
		panic(err2)
	}

	// Typed literal with URI
	result, err := cli.Query(
		ctx,
		"select * where "+
			"{ ?s <http://dbpedia.org/ontology/wikiPageLength> $1 . } LIMIT 1",
		Param{
			Ordinal: 1,
			Value: Literal{
				Value:    "76516",
				DataType: URI("http://www.w3.org/2001/XMLSchema#nonNegativeInteger"),
			},
		},
	)
	if err != nil {
		panic(err)
	}
	log.Println(result.Next())
	if err := result.Close(); err != nil {
		panic(err)
	}
}

func ExampleClient_Query_typed_literal_prefixed_name() {
	cli, err := New("http://ja.dbpedia.org/sparql",
		MaxIdleConns(100),
		IdleConnTimeout(90*time.Second),
		Timeout(30*time.Second),
		Prefix("dbpj", "http://ja.dbpedia.org/resource/"),
		Prefix("dbp-owl", "http://dbpedia.org/ontology/"),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	if err2 := cli.Ping(ctx); err2 != nil {
		panic(err2)
	}

	// Typed literal with PrefixedName
	result, err := cli.Query(
		ctx,
		`select * where { ?s <http://dbpedia.org/ontology/birthYear> $1 . } LIMIT 1`,
		Param{
			Ordinal: 1,
			Value: Literal{
				Value:    "1995",
				DataType: PrefixedName("xsd:gYear"),
			},
		},
	)
	if err != nil {
		panic(err)
	}
	log.Println(result.Next())
	if err := result.Close(); err != nil {
		panic(err)
	}
}

func ExampleClient_Query_parameterized_subject_and_object() {
	cli, err := New("http://ja.dbpedia.org/sparql",
		MaxIdleConns(100),
		IdleConnTimeout(90*time.Second),
		Timeout(30*time.Second),
		Prefix("dbpj", "http://ja.dbpedia.org/resource/"),
		Prefix("dbp-owl", "http://dbpedia.org/ontology/"),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	if err2 := cli.Ping(ctx); err2 != nil {
		panic(err2)
	}

	// Parameterized subject and object
	result, err := cli.Query(ctx, `select * where { $1 $2 ?o . } LIMIT 1000`,
		Param{
			Ordinal: 1,
			Value:   PrefixedName("dbpj:有安杏果"),
		},
		Param{
			Ordinal: 2,
			Value:   PrefixedName("dbp-owl:birthYear"),
		},
	)
	if err != nil {
		panic(err)
	}
	log.Println(result.Next())
	if err := result.Close(); err != nil {
		panic(err)
	}
}

func ExampleClient_Query_ask() {
	cli, err := New("http://ja.dbpedia.org/sparql",
		MaxIdleConns(100),
		IdleConnTimeout(90*time.Second),
		Timeout(30*time.Second),
		Prefix("dbpj", "http://ja.dbpedia.org/resource/"),
		Prefix("dbp-owl", "http://dbpedia.org/ontology/"),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	if err2 := cli.Ping(ctx); err2 != nil {
		panic(err2)
	}

	// Ask
	result, err := cli.Query(
		ctx,
		`ask { dbpj:有安杏果 dbp-owl:birthYear "1995"^^xsd:gYear . }`,
	)
	if err != nil {
		panic(err)
	}
	log.Println(result.Boolean())
	if err := result.Close(); err != nil {
		panic(err)
	}
}

func TestClient_Query(t *testing.T) {
	t.Run("request error", func(t *testing.T) {
		c, err := New("foo")
		if err != nil {
			t.Error(err)
			return
		}
		if _, err := c.Query(context.Background(), ""); err == nil {
			t.Errorf("Client.Query() error = %v", err)
		}
	})
	t.Run("not ok", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "", http.StatusBadRequest)
			},
		))

		c := &Client{
			HTTPClient: *server.Client(),
			Logger:     *logger.New(),
			Endpoint:   server.URL,
		}
		if _, err := c.Query(context.Background(), ""); err == nil {
			t.Errorf("Client.Query() error = %v", err)
			return
		}
	})
	t.Run("not json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = fmt.Fprint(w, "not json")
			},
		))

		c := &Client{
			HTTPClient: *server.Client(),
			Logger:     *logger.New(),
			Endpoint:   server.URL,
		}
		if _, err := c.Query(context.Background(), ""); err == nil {
			t.Errorf("Client.Query() error = %v", err)
			return
		}
	})
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = fmt.Fprint(w, `<?xml version="1.0"?>
<sparql xmlns="http://www.w3.org/2005/sparql-results#">
  <head>
    <variable name="x"/>
  </head>
  <results>
    <result> 
      <binding name="x"><bnode>r2</bnode></binding>
    </result>
  </results>
</sparql>`)
			},
		))

		c := &Client{
			HTTPClient: *server.Client(),
			Logger:     *logger.New(),
			Endpoint:   server.URL,
			prefixes:   map[string]URI{"foo": "bar"},
		}
		result, err := c.Query(context.Background(), "", Param{
			Ordinal: 0,
			Value:   1,
		})
		if err != nil {
			t.Errorf("Client.Query() error = %v", err)
			return
		}
		if got, want := result.Variables(), []string{"x"}; !reflect.DeepEqual(got, want) {
			t.Errorf("result.Variables() = %v, want %v", got, want)
		}
		bindings, err := result.Next()
		if err != nil {
			t.Errorf("iter.Next() error = %v", err)
			return
		}
		if want := map[string]Value{"x": BNode("r2")}; !reflect.DeepEqual(bindings, want) {
			t.Errorf("Client.Query() = %v, want %v", bindings, want)
		}
		if _, err = result.Next(); err != io.EOF {
			t.Errorf("iter.Next() error = %v", err)
			return
		}
	})
}

func BenchmarkClient_request(b *testing.B) {
	b.Run("query", func(b *testing.B) {
		client, err := New("endpoint")
		if err != nil {
			b.Fatal(err)
		}
		qs := make([]string, 0, b.N)
		for i := 0; i < b.N; i++ {
			qs = append(qs, fmt.Sprintf("query-%d", i))
		}

		ctx := context.Background()
		query := strings.Join(qs, ",")
		b.ResetTimer()
		if _, err := client.request(ctx, query); err != nil {
			b.Fatal(err)
		}
	})
	b.Run("params", func(b *testing.B) {
		client, err := New("endpoint")
		if err != nil {
			b.Fatal(err)
		}
		params := make([]Param, 0, b.N)
		for i := 0; i < b.N; i++ {
			params = append(params, Param{
				Ordinal: i,
				Value:   i,
			})
		}

		ctx := context.Background()
		b.ResetTimer()
		if _, err := client.request(ctx, "", params...); err != nil {
			b.Fatal(err)
		}
	})
}
