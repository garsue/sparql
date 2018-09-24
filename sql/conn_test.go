package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/garsue/sparql"
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
	db := sql.OpenDB(NewConnector(
		nil,
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

	//noinspection SqlDialectInspection
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

func TestRows_Columns(t *testing.T) {
	type fields struct {
		queryResult *sparql.QueryResult
		processed   int64
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "success",
			fields: fields{
				queryResult: &sparql.QueryResult{
					Head: sparql.Head{
						Vars: []string{"foo"},
					},
				},
			},
			want: []string{"foo"},
		},
		{
			name: "ASK",
			fields: fields{
				queryResult: &sparql.QueryResult{
					Head: sparql.Head{
						Vars: []string{},
					},
				},
			},
			want: []string{"boolean"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Rows{
				queryResult: tt.fields.queryResult,
				processed:   tt.fields.processed,
			}
			if got := r.Columns(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rows.Columns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRows_Close(t *testing.T) {
	r := &Rows{
		queryResult: &sparql.QueryResult{},
	}
	if err := r.Close(); err != nil {
		t.Errorf("Rows.Close() error = %v", err)
	}
	if r.queryResult != nil {
		t.Errorf("got nil want %v", r.queryResult)
	}
}

func TestRows_Next(t *testing.T) {
	t.Run("ASK", func(t *testing.T) {
		r := &Rows{
			queryResult: &sparql.QueryResult{
				Head:    sparql.Head{},
				Boolean: true,
			},
		}
		dest := make([]driver.Value, 1)
		if err := r.Next(dest); err != nil {
			t.Errorf("Rows.Next() error = %v", err)
		}
		want := []driver.Value{true}
		if !reflect.DeepEqual(dest, want) {
			t.Errorf("got %v want %v", dest, want)
		}
		if err := r.Next(dest); err != io.EOF {
			t.Errorf("Rows.Next() error = %v, wantErr %v", err, io.EOF)
		}
	})
	t.Run("success", func(t *testing.T) {
		r := &Rows{
			queryResult: &sparql.QueryResult{
				Head: sparql.Head{
					Vars: []string{"foo"},
				},
				Results: sparql.Results{
					Bindings: []sparql.Binding{
						{
							"foo": struct {
								Type     sparql.Type `json:"type"`
								DataType sparql.IRI  `json:"datatype"`
								XMLLang  string      `json:"xml:lang"`
								Value    interface{} `json:"value"`
							}{Value: 1},
						},
					},
				},
			},
		}
		dest := make([]driver.Value, 1)
		if err := r.Next(dest); err != nil {
			t.Errorf("Rows.Next() error = %v", err)
		}
		want := []driver.Value{1}
		if !reflect.DeepEqual(dest, want) {
			t.Errorf("got %v want %v", dest, want)
		}
		if err := r.Next(dest); err != io.EOF {
			t.Errorf("Rows.Next() error = %v, wantErr %v", err, io.EOF)
		}
	})
}

func TestConn_QueryContext(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		client, err := sparql.New("foo")
		if err != nil {
			t.Fatal(err)
		}
		c := &Conn{
			Client: client,
		}
		if _, err := c.QueryContext(context.Background(), "", nil); err == nil {
			t.Errorf("Conn.QueryContext() error = %v", err)
			return
		}
	})
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "{}")
			},
		))
		client, err := sparql.New(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		c := &Conn{
			Client: client,
		}
		got, err := c.QueryContext(context.Background(), "", []driver.NamedValue{
			{
				Name:    "foo",
				Ordinal: 0,
				Value:   1,
			},
		})
		if err != nil {
			t.Errorf("Conn.QueryContext() error = %v", err)
			return
		}
		want := driver.Rows(&Rows{
			queryResult: &sparql.QueryResult{
				Head:    sparql.Head{},
				Results: sparql.Results{},
			},
		})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Conn.QueryContext() = %+v, want %+v", got, want)
		}
	})
}

func TestConn_Ping(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "ok")
		},
	))
	client, err := sparql.New(server.URL)
	if err != nil {
		t.Fatal(err)
		return
	}
	c := &Conn{
		Client: client,
	}
	if err := c.Ping(context.Background()); err != nil {
		t.Errorf("Conn.Ping() error = %v", err)
	}
}

func TestConn_Prepare(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("got nil want error")
		}
	}()
	c := &Conn{
		Client: &sparql.Client{},
	}
	if _, err := c.Prepare(""); err != nil {
		t.Errorf("Conn.Prepare() error = %v", err)
		return
	}
}

func TestConn_Close(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "ok")
		},
	))
	client, err := sparql.New(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	c := &Conn{
		Client: client,
	}
	if err := c.Close(); err != nil {
		t.Errorf("Conn.Close() error = %v", err)
	}
}

func TestConn_Begin(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("got nil want error")
		}
	}()
	c := &Conn{
		Client: &sparql.Client{},
	}
	if _, err := c.Begin(); err != nil {
		t.Errorf("Conn.Begin() error = %v", err)
		return
	}
}
