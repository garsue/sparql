package sparql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/garsue/sparql/client"
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
		var p, o client.URI
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
		client.WithHTTPClient(&http.Client{
			Timeout: 5 * time.Second,
		}),
		client.WithPrefix("hj", "http://hojin-info.go.jp/ns/domain/biz/1#"),
		client.WithPrefix("ic", "http://imi.go.jp/ns/core/rdf#"),
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

type mockQueryResult struct {
	variables []string
	boolean   bool
	result    map[string]client.Value
	err       error
}

func (m *mockQueryResult) Variables() []string {
	return m.variables
}

func (m *mockQueryResult) Next() (map[string]client.Value, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

func (m *mockQueryResult) Boolean() (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.boolean, nil
}

func (m *mockQueryResult) Close() error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func TestRows_Columns(t *testing.T) {
	rows := Rows{
		queryResult: &mockQueryResult{
			variables: []string{"foo"},
		},
	}
	want := []string{"foo"}
	if got := rows.Columns(); !reflect.DeepEqual(got, want) {
		t.Errorf("Rows.Columns() = %v, want %v", got, want)
	}
}

func TestRows_Close(t *testing.T) {
	r := &Rows{
		queryResult: &mockQueryResult{},
	}
	if err := r.Close(); err != nil {
		t.Errorf("Rows.Close() error = %v", err)
	}
}

// nolint: gocyclo, dupl
func TestRows_Next(t *testing.T) {
	t.Run("ASK", func(t *testing.T) {
		r := &Rows{
			queryResult: &mockQueryResult{
				boolean: true,
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
	})
	t.Run("ASK error", func(t *testing.T) {
		askErr := errors.New("ASK error")
		r := &Rows{
			queryResult: &mockQueryResult{
				err: askErr,
			},
		}
		dest := make([]driver.Value, 1)
		if err := r.Next(dest); err != askErr {
			t.Errorf("Rows.Next() error = %v", err)
		}
	})
	t.Run("SELECT", func(t *testing.T) {
		r := &Rows{
			queryResult: &mockQueryResult{
				variables: []string{"foo"},
				result: map[string]client.Value{
					"foo": client.Literal{Value: "1"},
				},
			},
		}
		dest := make([]driver.Value, 1)
		if err := r.Next(dest); err != nil {
			t.Errorf("Rows.Next() error = %v", err)
		}
		want := []driver.Value{"1"}
		if !reflect.DeepEqual(dest, want) {
			t.Errorf("got %v want %v", dest, want)
		}
	})
	t.Run("SELECT error", func(t *testing.T) {
		selectErr := errors.New("SELECT error")
		r := &Rows{
			queryResult: &mockQueryResult{
				variables: []string{"foo"},
				err:       selectErr,
			},
		}
		dest := make([]driver.Value, 1)
		if err := r.Next(dest); err != selectErr {
			t.Errorf("Rows.Next() error = %v", err)
		}
	})
}

// nolint: scopelint
func Test_scan(t *testing.T) {
	type args struct {
		b client.Value
	}
	tests := []struct {
		name string
		args args
		want driver.Value
	}{
		{
			name: "uri",
			args: args{
				b: client.URI("http://www.w3.org/2001/XMLSchema#foo"),
			},
			want: client.URI("http://www.w3.org/2001/XMLSchema#foo"),
		},
		{
			name: "literal",
			args: args{
				b: client.Literal{
					DataType: client.URI("http://www.w3.org/2001/XMLSchema#integer"),
					Value:    "1",
				},
			},
			want: "1",
		},
		{
			name: "without timezone",
			args: args{
				b: client.Literal{
					DataType: client.URI("http://www.w3.org/2001/XMLSchema#dateTime"),
					Value:    "2015-11-19T00:10:11",
				},
			},
			want: time.Date(2015, time.November, 19, 0, 10, 11, 0,
				time.UTC,
			),
		},
		{
			name: "with timezone",
			args: args{
				b: client.Literal{
					DataType: client.URI("http://www.w3.org/2001/XMLSchema#dateTime"),
					Value:    "2015-11-19T09:10:11+09:00",
				},
			},
			want: time.Date(2015, time.November, 19, 9, 10, 11, 0,
				time.FixedZone("", 9*60*60),
			),
		},
		{
			name: "UTC",
			args: args{
				b: client.Literal{
					DataType: client.URI("http://www.w3.org/2001/XMLSchema#dateTime"),
					Value:    "2015-11-19T00:10:11Z",
				},
			},
			want: time.Date(2015, time.November, 19, 0, 10, 11, 0,
				time.UTC,
			),
		},
		{
			name: "nano",
			args: args{
				b: client.Literal{
					DataType: client.URI("http://www.w3.org/2001/XMLSchema#dateTime"),
					Value:    "2015-11-19T00:10:11.12345",
				},
			},
			want: time.Date(2015, time.November, 19, 0, 10, 11, 123450000,
				time.UTC,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := scan(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConn_QueryContext(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		cli, err := client.New("foo")
		if err != nil {
			t.Fatal(err)
		}
		c := &Conn{
			Client: cli,
		}
		if _, err := c.QueryContext(context.Background(), "", nil); err == nil {
			t.Errorf("Conn.QueryContext() error = %v", err)
			return
		}
	})
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = fmt.Fprint(w, `<sparql><head></head><results></results></sparql>`)
			},
		))
		cli, err := client.New(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		c := &Conn{
			Client: cli,
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
		if got, want := got.Columns(), []string{}; !reflect.DeepEqual(got, want) {
			t.Errorf("Conn.QueryContext() = %+v, want %+v", got, want)
		}
	})
}

func TestConn_Ping(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, "ok")
		},
	))
	cli, err := client.New(server.URL)
	if err != nil {
		t.Fatal(err)
		return
	}
	c := &Conn{
		Client: cli,
	}
	if err := c.Ping(context.Background()); err != nil {
		t.Errorf("Conn.Ping() error = %v", err)
	}
}

func TestConn_Prepare(t *testing.T) {
	c := &Conn{
		Client: &client.Client{},
	}
	if _, err := c.Prepare(""); err != nil {
		t.Errorf("Conn.Prepare() error = %v", err)
		return
	}
}

func TestConn_Close(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprint(w, "ok")
		},
	))
	cli, err := client.New(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	c := &Conn{
		Client: cli,
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
		Client: &client.Client{},
	}
	if _, err := c.Begin(); err != nil {
		t.Errorf("Conn.Begin() error = %v", err)
		return
	}
}
