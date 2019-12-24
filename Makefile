TAGS ?= ""

tidy:
	go mod tidy

test:
	go test -cover -tags ${TAGS} ./...
	make tidy

