build:
	@go build -o bin/vgen src/*.go
	@echo "binary created in bin/vgen"

compile_simple:
	@make build && ./bin/vgen examples/simple/simple.go
run_simple:
	@go run examples/simple/*.go
all_simple:
	@make compile_simple
	@make run_simple

compile_list:
	@make build && ./bin/vgen examples/list/list.go
run_list:
	@go run examples/list/*.go
all_list:
	@make compile_list
	@make run_list

compile_types:
	@make build && ./bin/vgen examples/types/types.go
run_types:
	@go run examples/types/*.go
all_types:
	@make compile_types
	@make run_types
