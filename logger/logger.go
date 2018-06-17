package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

// LogLevel stands for log levels
type LogLevel int

// Log levels
const (
	None LogLevel = iota
	Debug
	Error
)

// Logger offers logging methods
type Logger struct {
	Debug *log.Logger
	Error *log.Logger
}

// New returns a logger
func New() *Logger {
	switch env {
	case None:
		return &Logger{
			Debug: log.New(ioutil.Discard, "SPARQL:DEBUG:", log.LstdFlags),
			Error: log.New(ioutil.Discard, "SPARQL:ERROR:", log.LstdFlags),
		}
	case Error:
		return &Logger{
			Debug: log.New(ioutil.Discard, "SPARQL:DEBUG:", log.LstdFlags),
			Error: log.New(os.Stderr, "SPARQL:ERROR:", log.LstdFlags),
		}
	case Debug:
		return &Logger{
			Debug: log.New(os.Stdout, "SPARQL:DEBUG:", log.LstdFlags),
			Error: log.New(os.Stderr, "SPARQL:ERROR:", log.LstdFlags),
		}
	}
	log.Panicf("unknown log level %d", env)
	return nil
}

// LogCloseError reports a closing error
func (l *Logger) LogCloseError(closer io.Closer) {
	if err := closer.Close(); err != nil {
		l.Error.Println(err)
	}
}
