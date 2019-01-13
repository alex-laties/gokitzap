# gokit-zap [![Build Status](https://travis-ci.com/alex-laties/gokitzap.svg?branch=master)](https://travis-ci.com/alex-laties/gokitzap)

A quick adapter to allow the use of [`zap`](https://github.com/uber-go/zap) as the underlying logger for [`go-kit/log`](https://github.com/go-kit/kit/tree/master/log).
Transparently translates log levels in go-kit to zap log levels.

# Usage

```
package main
import (
  "time"

  "github.com/alex-laties/gokitzap"
  kitlog "github.com/go-kit/kit/log"
  "github.com/go-kit/kit/log/level"
  "go.uber.org/zap"
)

var logger kitlog.Logger

func main() {
  mainStart := time.Now()
  zl := zap.NewDevelopment()

  logger = gokitzap.FromZLogger(zl)

  logger.Log("message", "hello world")
  level.Debug(logger).Log("message", "levels work too")

  level.Info("message", "startup time", "t", time.Since(mainStart))
}
```

## Uhhh... why?

`go-kit/log` is a reasonable logging option, but can lack performance in certain situations.
It's not unreasonable to want more performance from your logging framework, but it can be difficult to transition everything to a completely different library like `zap` in one pass.

`gokit-zap` allows you to adopt `zap` under the hood while still offering the `go-kit/log` interface, allowing for immediate performance benefits while one transitions to `zap` completely.

## Benchmarks
```
# make benchmark
PASS
ok  	github.com/alex-laties/gokitzap	0.011s
goos: darwin
goarch: amd64
pkg: github.com/alex-laties/gokitzap/benchmarks
BenchmarkGoKit-6            	 1000000	      1114 ns/op	    1304 B/op	      22 allocs/op
BenchmarkGoKitLevels-6      	 1000000	      1428 ns/op	    1770 B/op	      28 allocs/op
BenchmarkZapSugar-6         	 3000000	       484 ns/op	     163 B/op	       1 allocs/op
BenchmarkZapSugarLevels-6   	 3000000	       536 ns/op	     163 B/op	       1 allocs/op
BenchmarkGKZ-6              	 2000000	       610 ns/op	     705 B/op	       1 allocs/op
BenchmarkGKZLevels-6        	 2000000	       807 ns/op	    1026 B/op	       4 allocs/op
PASS
ok  	github.com/alex-laties/gokitzap/benchmarks	11.075s
```
