package gokitzap

import (
	"errors"

	"github.com/go-kit/kit/log/level"
	"go.uber.org/zap"
)

var ErrNilLogger = errors.New("provided logger is nil")
var ErrUnmatchedKeyVals = errors.New("got unmatched keys/values")

type GKZLogger struct {
	s *zap.SugaredLogger
	messageKey string
}

func FromZLogger(l *zap.Logger) (*GKZLogger, error) {
	if l == nil {
		return nil, ErrNilLogger
	}
	return FromZSLogger(l.Sugar())
}

func FromZSLogger(s *zap.SugaredLogger) (*GKZLogger, error) {
	if s == nil {
		return nil, ErrNilLogger
	}
	return &GKZLogger{
		s: s,
		messageKey: "message",
	}, nil
}

func (gkz *GKZLogger) SetMessageKey(key string) {
	gkz.messageKey = key
}

func (gkz *GKZLogger) Log(keyvals ...interface{}) error {
	if len(keyvals) == 0 {
		return nil
	}

	// extract message key and value
	var msg string
	for i := 0; i < len(keyvals); i += 2 {
		if keyvals[i] == gkz.messageKey {
			msg = keyvals[i+1].(string)

			if len(keyvals) == 2 {
				keyvals = nil
				break
			}

			switch {
			case i == 0: // head
				keyvals = keyvals[2:]
			case i+2 > len(keyvals) - 1: // tail
				keyvals = keyvals[0:i]
			default:
				keyvals = append(keyvals[0:i], keyvals[i+2:]...)
			}
			break
		}
	}

	// extract log level key and value
	l := level.InfoValue()
	for i := 0; i < len(keyvals); i += 2 {
		if keyvals[i] == level.Key() {
			

			l = keyvals[i+1].(level.Value)

			// have to strip log level from the result
			if len(keyvals) == 2 { // only a level is logged?
				keyvals = nil
				break
			}

			switch {
			case i == 0: // head
				keyvals = keyvals[2:]
			case i + 2 > len(keyvals) - 1: // tail
				keyvals  = keyvals[0:i]
			default: // middle 
				keyvals = append(keyvals[0:i], keyvals[i+2:]...)
			}
			break
		}
	}


	switch l {
	case level.DebugValue():
		gkz.s.Debugw(msg, keyvals...)
	case level.ErrorValue():
		gkz.s.Errorw(msg, keyvals...)
	case level.WarnValue():
		gkz.s.Warnw(msg, keyvals...)
	case level.InfoValue():
		gkz.s.Infow(msg, keyvals...)
	}

	return nil
}
