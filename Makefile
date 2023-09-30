build:
	@go build -o bin/vgen src/*.go
	@echo "binary created in bin/vgen"

example_test:
	@go run examples/test/*.go

example_test2:
	@go run examples/test2/*.go


