default:
	go run main.go --help

generate:
	go generate ./...

build: generate
	go build -o kb main.go

install: generate
	go install
