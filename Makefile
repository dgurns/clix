default: run

run:
	go run ./cmd/clix/main.go

bin:
	go build -o ./bin/clix ./cmd/clix/main.go

install:
	go install ./cmd/clix

reset:
	-rm ~/.clix/.env
	-rm ./bin/clix