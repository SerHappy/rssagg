build:
	go build -o build/server cmd/main.go

run:
	go run cmd/main.go

start: build
	./build/server

install:
	go mod tidy
