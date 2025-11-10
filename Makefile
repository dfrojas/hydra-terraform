mod:
	go mod init hydratf

sum:
	go mod tidy

build:
	go build -o hydratf

init:
	./hydratf init --source example-module/main.tf --name mocked

generate:
	./hydratf generate
