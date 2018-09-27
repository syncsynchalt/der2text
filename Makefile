all: read

read:
	go build ./cmds/der2text
	go build ./cmds/text2der

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
