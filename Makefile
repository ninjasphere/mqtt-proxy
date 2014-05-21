
all:
	scripts/build.sh

clean:
	rm bin/mqtt-proxy || true
	rm -rf .gopath || true

test:
	go test ./...

.PHONY: all	clean test
