build:
	@go build -o bin/vgen src/*.go

test:
	@go run examples/test/run/main.go
