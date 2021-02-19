.env:
	cp .env.example .env

.PHONY: deps
deps: .env
	docker-compose run --rm golang go mod tidy

.PHONY: test
test: .env
	docker-compose run --rm -e GOOS=linux golang go test -v ./...

.PHONY: build
build: .env
	rm -rf bin
	docker-compose run --rm golang go build -ldflags="-s -w" -o bin/marvel-characters-api cmd/api/main.go

.PHONY: start
start: .env
	docker-compose run --rm -p 8080:8080 -e GOOS=linux golang go run cmd/api/main.go

.PHONY: fmt
fmt: .env
	docker-compose run --rm golang go fmt ./...

.PHONY: genMocks
genMocks: .env
	docker-compose run --rm mockery --all

.PHONY: genSwagger
genSwagger: .env
	docker-compose run --rm -e GOOS=linux golang sh scripts/gen-swagger.sh
