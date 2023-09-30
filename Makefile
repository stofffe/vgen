build:
	@go build -o bin/vgen src/*.go
	@echo "binary created in bin/vgen"

example_test1:
	@go run examples/test1/*.go

example_test2:
	@go run examples/test2/*.go

example_test3:
	@go run examples/test3/*.go

example_test4:
	@go run examples/test4/*.go

all_test1:
	@make build && ./bin/vgen examples/test1/test1.go && make example_test1

all_test2:
	@make build && ./bin/vgen examples/test2/test2.go && make example_test2

all_test3:
	@make build && ./bin/vgen examples/test3/test3.go && make example_test3

all_test4:
	@make build && ./bin/vgen examples/test4/test4.go && make example_test4
