TAGS ?= ""

test:
	go test -cover -race -tags ${TAGS} ./...
