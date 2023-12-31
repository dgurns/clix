default: run

run:
	go run ./cmd/clix/main.go

build:
	go build -o ./bin/clix ./cmd/clix/main.go

install:
	go install ./cmd/clix
