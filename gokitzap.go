package gokitzap

import (
	"errors"

	"github.com/go-kit/kit/log/level"
	"go.uber.org/zap"
)

var ErrNilLogger = errors.New("provided logger is nil")
var ErrUnmatchedKeyVals = errors.New("got unmatched keys/values")
var ErrOnlyLevel = errors.New("only got a level in the log")

type GKZLogger struct {
	s *zap.SugaredLogger
}

func FromZLogger(l *zap.Logger) (*GKZLogger, error) {
	if l == nil {
		return nil, ErrNilLogger
	}
	return &GKZLogger{
		s: l.Sugar(),
	}, nil
}

func FromZSLogger(s *zap.SugaredLogger) (*GKZLogger, error) {
	if s == nil {
		return nil, ErrNilLogger
	}
	return &GKZLogger{
		s: s,
	}, nil
}

func (gkz *GKZLogger) Log(keyvals ...interface{}) error {
	if len(keyvals) == 0 {
		return nil
	}
	// extract log level key and value
	var l level.Value
	for i := 0; i < len(keyvals); i += 2 {
		if keyvals[i] == level.Key() {
			if len(keyvals) == 2 { // only a level is logged?
				return ErrOnlyLevel
			}

			l = keyvals[i+1].(level.Value)

			// have to strip log level from the result
			switch {
			case i == 0: // head
				keyvals = keyvals[2:]
			case i + 2 > len(keyvals) - 1: // tail
				keyvals  = keyvals[0:i]
			default: // middle 
				keyvals = append(keyvals[0:i], keyvals[i+2:])
			}
			break
		}
	}

	switch l {
	case level.DebugValue():
		gkz.s.Debug(keyvals)
	case level.ErrorValue():
		gkz.s.Error(keyvals)
	case level.WarnValue():
		gkz.s.Warn(keyvals)
	case level.InfoValue():
		gkz.s.Info(keyvals)
	default: // no level so go for Info
		gkz.s.Info(keyvals)
	}

	return nil
}
