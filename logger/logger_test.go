package logger

import (
	"errors"
	"log"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		if got := New(); got == nil {
			t.Errorf("New() = nil")
		}
	})
}

type closerFunc func() error

func (e closerFunc) Close() error {
	return e()
}

func ExampleLogger_LogCloseError() {
	logger := &Logger{
		Error: log.New(os.Stdout, "", 0),
	}
	logger.LogCloseError(closerFunc(func() error { return errors.New("error") }))
	logger = &Logger{
		Error: log.New(os.Stdout, "", 0),
	}
	logger.LogCloseError(closerFunc(func() error { return nil }))

	// Output: error
}
