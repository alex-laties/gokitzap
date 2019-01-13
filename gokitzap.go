package gokitzap

import (
	"errors"

	"github.com/go-kit/kit/log/level"
	"go.uber.org/zap"
)

// ErrNilLogger is returned when a nil logger is passed to an adapter function like `FromZLogger'
var ErrNilLogger = errors.New("provided logger is nil")
// ErrUnmatchedKeyVals is returned when a set of keyval pairs is odd, i.e. a key is missing a value.
var ErrUnmatchedKeyVals = errors.New("got unmatched keys/values")

// GKZLogger is a wrapper for *zap.Logger that implements the github.com/go-kit/kit/log.Logger interface
type GKZLogger struct {
	s *zap.SugaredLogger
	messageKey string
}

// FromZLogger takes a valid *zap.Logger and wraps it in a *GKZLogger
// Call `SetMessageKey` if the message key differs from the default of "message"
func FromZLogger(l *zap.Logger) (*GKZLogger, error) {
	if l == nil {
		return nil, ErrNilLogger
	}
	return FromZSLogger(l.Sugar())
}

// FromZSLogger takes a valid *zap.SugaredLogger and wraps it in a *GKZLogger
// Call `SetMessageKey` if the message key differs from the default of "message"
func FromZSLogger(s *zap.SugaredLogger) (*GKZLogger, error) {
	if s == nil {
		return nil, ErrNilLogger
	}
	return &GKZLogger{
		s: s,
		messageKey: "message",
	}, nil
}

// SetMessageKey changes the field key that the GKZLogger will extract the message from
// By default, the message key is "message"
func (gkz *GKZLogger) SetMessageKey(key string) {
	gkz.messageKey = key
}

// Log translates a kitlog set of keyval pairs to a leveled *zap.SugaredLogger.(Debugw, Errorw, Infow, etc) call.
// Log levels are determined using values from github.com/go-kit/kit/log/level constants.
// If no level is determined, Log defaults to logging at INFO level.
// Messages are determined using the message key set on the GKZLogger.
// If no message is determined, the message is left blank.
func (gkz *GKZLogger) Log(keyvals ...interface{}) error {
	if len(keyvals) == 0 {
		return nil
	}

	if len(keyvals) % 2 != 0 {
		return ErrUnmatchedKeyVals
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
