PKGNAME    ?= github.com/nikonov1101/kb
VERSION    ?= $(shell git describe --long --always --dirty --broken --abbrev=8 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date --utc --iso-8601=minutes)
LDFLAGS    := "-X $(PKGNAME)/version.Commit=$(VERSION) -X $(PKGNAME)/version.BuildTime=$(BUILD_DATE)"

default:
	go run main.go --help

build:
	@go build \
		-buildvcs=false \
		-trimpath \
		-ldflags $(LDFLAGS) \
		-o kb main.go

install:
	go install

