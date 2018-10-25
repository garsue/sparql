package sparql

import (
	"testing"
)

func TestDriver_Open(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		d := &Driver{}
		got, err := d.Open("name")
		if err != nil {
			t.Errorf("Driver.Open() error = %v", err)
			return
		}
		if got == nil {
			t.Errorf("Driver.Open() = nil")
		}
	})
}

func TestDriver_OpenConnector(t *testing.T) {
	d := &Driver{}
	got, err := d.OpenConnector("name")
	if err != nil {
		t.Errorf("Driver.OpenConnector() error = %v", err)
		return
	}
	if got == nil {
		t.Errorf("Driver.OpenConnector() = nil")
	}
}
