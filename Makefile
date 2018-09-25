all: read

read:
	go build github.com/syncsynchalt/der2text/cmds/der2text

test:
	@set -e; for d in $$(dirname $$(find . -name '*_test.go')); do echo testing in $$d; go test $$d; done

clean:
	rm -f der2text

.PHONY: read clean test
