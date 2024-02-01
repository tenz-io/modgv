
.PHONY build:
build:
	@echo "Building the project..."
	@go build -o bin/modgv main.go


.PHONY test:
test: export MODGV_DST_NODE := golang.org/x/tools
test:
	@echo "Running tests..."
	@go mod graph | go run main.go


.PHONY test-png:
test-png: export MODGV_DST_NODE := golang.org/x/tools
test-png:
	@echo "Running tests..."
	@go mod graph | go run main.go | dot -Grankdir=TB -Tpng -o testdata/test.png
	@open testdata/test.png