TAGS ?= ""

install: packr
	cd ./genny && go install -tags ${TAGS} -v .
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
