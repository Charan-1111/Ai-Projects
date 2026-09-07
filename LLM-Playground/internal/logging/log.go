package logging

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
)

type Log struct {
	Log   zerolog.Logger
	Close func()
	once  sync.Once
}

func (log *Log) Initialize() {
	log.once = sync.Once{}
	log.Log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	log.Close = func() {}
}

func (log *Log) WithRequestID(requestID string) zerolog.Logger {
	if log == nil {
		return zerolog.Nop()
	}
	if requestID == "" {
		return log.Log
	}
	return log.Log.With().Str("requestId", requestID).Logger()
}

func (log *Log) CloseFunc() {
	log.once.Do(func() {
		if log.Close != nil {
			log.Close()
		}
	})
}
