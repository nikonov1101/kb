default:
	go run main.go --help

build: 
	go generate ./...
	go build -o kb main.go

install:
	go generate ./...
	go install

