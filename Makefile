default:
	go run main.go --help

build:
	go build -o kb main.go

install:
	go install
