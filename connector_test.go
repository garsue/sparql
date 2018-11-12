package sparql

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/garsue/sparql/client"
)

func TestNewConnector(t *testing.T) {
	opts := []client.Option{client.Timeout(30 * time.Second)}
	want := &Connector{
		driver:  nil,
		Name:    "name",
		options: opts,
	}
	if got := NewConnector(nil, "name", opts...); !reflect.DeepEqual(got, want) {
		t.Errorf("NewConnector() = %+v, want %+v", got, want)
	}
}

func TestConnector_Connect(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		c := &Connector{
			driver: nil,
			Name:   "name",
			options: []client.Option{
				func(*client.Client) error {
					return errors.New("error")
				},
			},
		}
		if _, err := c.Connect(context.Background()); err == nil {
			t.Errorf("Connector.Connect() error = %v", err)
			return
		}
	})
	t.Run("success", func(t *testing.T) {
		c := &Connector{
			driver:  nil,
			Name:    "name",
			options: nil,
		}
		got, err := c.Connect(context.Background())
		if err != nil {
			t.Errorf("Connector.Connect() error = %v", err)
			return
		}
		if got, ok := got.(*Conn); ok && got.Client == nil {
			t.Errorf("Connector.Connect() = %v", got)
		}
	})
}

func TestConnector_Driver(t *testing.T) {
	c := &Connector{
		driver:  &Driver{},
		Name:    "name",
		options: nil,
	}
	if got := c.Driver(); !reflect.DeepEqual(got, c.driver) {
		t.Errorf("Connector.Driver() = %v, want %v", got, c.driver)
	}
}
