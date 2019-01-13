package gokitzap

import (
	"testing"

	"github.com/go-kit/kit/log/level"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func message(m string, fields ...zapcore.Field) observer.LoggedEntry {
	if fields == nil {
		fields = []zap.Field{}
	}
	return observer.LoggedEntry{Entry: zapcore.Entry{Message: m}, Context: fields}
}

func leveledMessage(m string, lvl zapcore.Level, fields ...zapcore.Field) observer.LoggedEntry {
	if fields == nil {
		fields = []zap.Field{}
	}
	return observer.LoggedEntry{Entry: zapcore.Entry{Message: m, Level: lvl}, Context: fields}
}

func TestBasic(t *testing.T) {
	core, observed := observer.New(zap.DebugLevel)
	l := zap.New(core)
	gkzl, _ := FromZLogger(l)
	assert.NoError(t, gkzl.Log("message", "hello world"))
	assert.NoError(t, gkzl.Log("message", "oh hai"))
	assert.Equal(
		t,
		[]observer.LoggedEntry{
			message("hello world"),
			message("oh hai"),
		}, 
		observed.AllUntimed(),
		"entry mismatch",
	)

	assert.Equal(t, ErrUnmatchedKeyVals, gkzl.Log("this should not work"))
	assert.Equal(t, ErrUnmatchedKeyVals, gkzl.Log("1", "2", "3"))
	assert.NoError(t, gkzl.Log())
}

func TestComplex(t *testing.T) {
	core, observed := observer.New(zap.DebugLevel)
	l := zap.New(core)
	gkzl, _ := FromZLogger(l)
	assert.NoError(t, gkzl.Log("message", "hello", "foo", 22, "bar", "wha"))
	assert.Equal(
		t,
		[]observer.LoggedEntry{
			message("hello", zap.Int("foo", 22), zap.String("bar", "wha")),
		},
		observed.AllUntimed(),
		"entry mismatch",
	)
}

func TestLevels(t *testing.T) {
	core, observed := observer.New(zap.DebugLevel)
	l := zap.New(core)
	gkzl, _ := FromZLogger(l)
	ll := level.NewFilter(gkzl, level.AllowAll())
	assert.NoError(t, level.Debug(ll).Log("message", "woah"))
	assert.NoError(t, level.Warn(ll).Log("message", "woah"))
	assert.NoError(t, level.Error(ll).Log("message", "woah"))
	assert.NoError(t, level.Info(ll).Log("message", "woah"))

	assert.Equal(
		t,
		[]observer.LoggedEntry{
			leveledMessage("woah", zap.DebugLevel),
			leveledMessage("woah", zap.WarnLevel),
			leveledMessage("woah", zap.ErrorLevel),
			leveledMessage("woah", zap.InfoLevel),
		},
		observed.AllUntimed(),
		"entry mismatch",
	)
}
