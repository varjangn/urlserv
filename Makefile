build:
	@go build -o bin/urlserv

run: build
	@./bin/urlserv

test:
	@go test -f ./ ...
