all: read

read:
	go build ./cmds/...

test:
	go test ./read/...
	go test ./write/...

clean:
	rm -f der2text text2der
	go clean -cache

vet:
	go vet --shadow ./...

fmt:
	go fmt ./...

.PHONY: read clean test
