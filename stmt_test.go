package sparql

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garsue/sparql/client"
)

//noinspection ALL
func TestStmt_QueryContext(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		db := sql.OpenDB(NewConnector("foo"))
		defer db.Close()
		stmt, err := db.Prepare("")
		if err != nil {
			t.Errorf("Conn.Prepare() error = %v", err)
			return
		}
		if _, err := stmt.QueryContext(context.Background()); err == nil {
			t.Errorf("Conn.QueryContext() error = %v", err)
			return
		}
	})
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, h *http.Request) {
			_, _ = w.Write([]byte(`<sparql><head></head><results></results></sparql>`))
		}))
		db := sql.OpenDB(NewConnector(server.URL))
		defer db.Close()
		stmt, err := db.Prepare("")
		if err != nil {
			t.Errorf("Conn.Prepare() error = %v", err)
			return
		}
		if _, err := stmt.QueryContext(context.Background()); err != nil {
			t.Errorf("Conn.QueryContext() error = %v", err)
			return
		}
	})
}

func TestStmt_Close(t *testing.T) {
	var s Stmt
	if err := s.Close(); err != nil {
		t.Errorf("Stmt.Close() error = %v", err)
	}
}

func TestStmt_NumInput(t *testing.T) {
	s := Stmt{
		Statement: &client.Statement{},
	}
	if got := s.NumInput(); got != -1 {
		t.Errorf("Stmt.NumInput() = %v, want -1", got)
	}
}

func TestStmt_Exec(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("Expect panic")
		}
	}()
	var s Stmt
	_, _ = s.Exec(nil)
}

func TestStmt_Query(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("Expect panic")
		}
	}()
	var s Stmt
	_, _ = s.Query(nil)
}
