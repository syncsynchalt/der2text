all: read

read:
	go build github.com/syncsynchalt/der2text/cmds/der2text

test:
	go test github.com/syncsynchalt/der2text/read/...

clean:
	rm -f der2text
	go clean -cache

.PHONY: read clean test
