package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWithHTTPClient(t *testing.T) {
	timeout := 30 * time.Second
	httpClient := &http.Client{
		Timeout: timeout,
	}
	client := Client{}
	if err := WithHTTPClient(httpClient)(&client); err != nil {
		t.Error(err)
		return
	}
	if got, want := client.HTTPClient.Timeout, timeout; got != want {
		t.Errorf("TestWithHTTPClient() = %v, want %v", got, want)
	}
}

func TestWithPrefix(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		prefix := "dbpj"
		uri := URI("http://ja.dbpedia.org/resource/")
		client := Client{
			prefixes: map[string]URI{},
		}
		if err := WithPrefix(prefix, uri)(&client); err != nil {
			t.Error(err)
			return
		}
		if got, want := client.prefixes[prefix], uri; got != want {
			t.Errorf("WithPrefix() = %v, want %v", got, want)
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		endpoint := "http://ja.dbpedia.org/sparql"
		got, err := New(endpoint)
		if err != nil {
			t.Errorf("New() error = %v", err)
			return
		}
		if got.Endpoint != endpoint {
			t.Errorf("New() = %s, want %s", got.Endpoint, endpoint)
		}
	})
	t.Run("error", func(t *testing.T) {
		endpoint := "http://ja.dbpedia.org/sparql"
		_, err := New(endpoint, func(*Client) error {
			return errors.New("error")
		})
		if err == nil {
			t.Errorf("New() error = %v", err)
			return
		}
	})
}

func TestClient_Close(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := &Client{
			HTTPClient: http.Client{
				Transport: http.DefaultTransport,
			},
		}
		if err := c.Close(); err != nil {
			t.Errorf("Client.Close() error = %v", err)
		}
	})
}

func TestClient_Ping(t *testing.T) {
	t.Run("request error", func(t *testing.T) {
		c, err := New("foo")
		if err != nil {
			t.Error(err)
			return
		}
		if err := c.Ping(context.Background()); err == nil {
			t.Errorf("Client.Ping() error = %v", err)
		}
	})
	t.Run("not ok", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "", http.StatusBadRequest)
			}),
		)
		c, err := New(server.URL)
		if err != nil {
			t.Error(err)
			return
		}
		if err := c.Ping(context.Background()); err == nil {
			t.Errorf("Client.Ping() error = %v", err)
		}
	})
	t.Run("ok", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = fmt.Fprint(w, "ok")
			}),
		)
		c, err := New(server.URL)
		if err != nil {
			t.Error(err)
			return
		}
		if err := c.Ping(context.Background()); err != nil {
			t.Errorf("Client.Ping() error = %v", err)
		}
	})
}
