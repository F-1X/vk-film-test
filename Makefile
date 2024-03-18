docker-up:
	docker-compose up

docker-down:
	docker-compose down

go-build:
	@go build -o bin/main ./cmd

go-run: go-build
	@./bin/main

go-install-deps:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go get -u gopkg.in/yaml.v2


coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out
	@rm coverage.out

#export PATH="$PATH:$(go env GOPATH)/bin"
# gen:
# 	mockgen -source=database/database.go

