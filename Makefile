TAGS ?= ""
GO_BIN ?= "go"

install: packr
	go install -tags ${TAGS} -v ./genny
	make tidy

tidy:
	go mod tidy

build: packr
	go build -v .
	make tidy

test: packr
	go test -cover -tags ${TAGS} ./...
	make tidy

packr:
	packr2
	make tidy
