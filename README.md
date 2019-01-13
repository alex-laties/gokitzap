# gokit-zap [![Build Status][ci-img]][ci]

A quick adapter to allow the use of [`zap`](https://github.com/uber-go/zap) as the underlying logger for [`go-kit/log`](https://github.com/go-kit/kit/tree/master/log).
Transparently translates log levels in go-kit to zap log levels.

## Uhhh... why?

`go-kit/log` is a reasonable logging option, but can lack performance in certain situations.
It's not unreasonable to want more performance from your logging framework, but it can be difficult to transition everything to a completely different library like `zap` in one pass.

`gokit-zap` allows you to adopt `zap` under the hood while still offering the `go-kit/log` interface, allowing for immediate performance benefits while one transitions to `zap` completely.
