.PHONY: test
test:
	@go test ./...

.PHONY: benchmarklong
benchmarklong:
	@go test -bench . -benchtime 1m ./...

.PHONY: benchmark
benchmark:
	@go test -bench . ./...

.DEFAULT_GOAL: test
