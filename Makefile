all:
	go build ./cmds/der2text
	go build ./cmds/text2der

deps:
	go get golang.org/x/text/encoding/unicode

test:
	go test ./read/...
	go test ./write/...
	@result=0; \
	for i in $$(ls samples/* | grep -v README); do \
		echo "=== $$i"; \
		./der2text $$i | ./text2der > /tmp/t.$$; \
		cmp $$i /tmp/t.$$; \
		if [ "$$?" != "0" ]; then \
			echo "Failure in file $$i" 1>&2; \
			result=1; \
		fi; \
		rm -f /tmp/t.$$; \
	done; \
	exit $$result

clean:
	rm -f der2text text2der
	go clean -cache

vet:
	go vet --shadow ./...

fmt:
	go fmt ./...

.PHONY: read clean test
