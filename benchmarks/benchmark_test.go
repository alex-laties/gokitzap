package benchmarks

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
	 
	"github.com/alex-laties/gokitzap"
	gokit "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	randomdata "github.com/Pallinder/go-randomdata"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ValType int

const (
	INT ValType = iota
	BOOL
	FLOAT
	STRING
	ERR
)

type LogLevelType int

const (
	DEBUG LogLevelType = iota
	INFO
	WARN
	ERROR
)

var (
	LOGLEVELS = []LogLevelType{DEBUG, INFO, WARN, ERROR}
	TYPES = []ValType{INT, BOOL, FLOAT, STRING, ERR}
	randomMessages [][]interface{}
)

func randomVal(t ValType) interface{} {
	switch t {
	case INT:
		return rand.Int()
	case BOOL:
		return rand.Int() % 2 == 0
	case FLOAT:
		return rand.Float64()
	case STRING:
		return fmt.Sprintf("%s_%s", randomdata.Adjective(), randomdata.Noun())
	case ERR:
		return fmt.Errorf("%s_%s", randomdata.Adjective(), randomdata.Noun())
	}
	return nil
}

// generate a random semi-structured key-val log line, between 1 and 10 fields + a message
func makeMessage() []interface{} {
	message := make([]interface{}, 0, 22)
	// always include a message
	message = append(message, "message")
	message = append(message, randomVal(STRING))
	keyvalPairs := rand.Intn(10) + 1 // rand.Intn(10) returns 0-9, but we want 1-10, hence +1
	for i := 0; i < keyvalPairs; i++ {
		// key
		message = append(message, randomVal(STRING))
		// val
		valtype := TYPES[rand.Intn(len(TYPES))]
		message = append(message, randomVal(valtype))
	}

	return message
}

func randomMessage() []interface{} {
	return randomMessages[rand.Intn(len(randomMessages))]
}

func randomLogLevel() LogLevelType {
	return LOGLEVELS[rand.Intn(len(LOGLEVELS))]
}

func makeZapLogger() *zap.Logger {
	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(
				zap.NewProductionEncoderConfig(),
			),
			zapcore.AddSync(ioutil.Discard),
			zap.DebugLevel,
		),
	)
}

func init() {
	rand.Seed(time.Now().UnixNano())
	randomMessages = make([][]interface{}, 0, 9999)

	for i := 0; i < 9999; i++ {
		randomMessages = append(randomMessages, makeMessage())
	}
}

func BenchmarkGoKit(b *testing.B) {
	l := gokit.With(gokit.NewJSONLogger(gokit.NewSyncWriter(ioutil.Discard)))
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Log(randomMessage())
		}
	})
}

func BenchmarkGoKitLevels(b *testing.B) {
	l := gokit.With(gokit.NewJSONLogger(gokit.NewSyncWriter(ioutil.Discard)))
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			switch randomLogLevel() {
			case DEBUG:
				level.Debug(l).Log(randomMessage())
			case ERROR:
				level.Error(l).Log(randomMessage())
			case WARN:
				level.Warn(l).Log(randomMessage())
			case INFO:
				level.Info(l).Log(randomMessage())
			}
		}
	})
}

func BenchmarkZapSugar(b *testing.B) {
	l := makeZapLogger().Sugar()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info(randomMessage())
		}
	})
}

func BenchmarkZapSugarLevels(b *testing.B) {
	l := makeZapLogger().Sugar()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			switch randomLogLevel() {
			case DEBUG:
				l.Debug(randomMessage())
			case WARN:
				l.Warn(randomMessage())
			case ERROR:
				l.Error(randomMessage())
			case INFO:
				l.Info(randomMessage()) 
			}
		}
	})
}

func BenchmarkGKZ(b *testing.B) {
	l, _ := gokitzap.FromZLogger(makeZapLogger())
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Log(randomMessage())
		}
	})
}

func BenchmarkGKZLevels(b *testing.B) {
	l, _ := gokitzap.FromZLogger(makeZapLogger())
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			switch randomLogLevel() {
			case DEBUG:
				level.Debug(l).Log(randomMessage())
			case WARN:
				level.Warn(l).Log(randomMessage())
			case ERROR:
				level.Error(l).Log(randomMessage())
			case INFO:
				level.Info(l).Log(randomMessage())
			}
		}
	})
}
