mod:
	go mod init hydratf

sum:
	go mod tidy

build:
	go build -o bin/hydratf ./cmd/hydratf

install:
	go install ./cmd/hydratf

run: build
	./bin/hydratf init --source test-data/main.tf --name localstack
	./bin/hydratf generate
