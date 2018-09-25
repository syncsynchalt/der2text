all: read

read:
	go build github.com/syncsynchalt/der2text/cmds/der2text

test:
	go test github.com/syncsynchalt/der2text/...

clean:
	rm -f der2text

.PHONY: read clean test
