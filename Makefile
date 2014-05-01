
all:
	scripts/build.sh

clean:
	rm bin/mqtt-proxy || true
	rm -rf .gopath || true

.PHONY: all	clean
