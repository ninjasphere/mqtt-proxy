
all: clean
	$(shell cd content; $(GOPATH)/bin/rice embed)
	$(shell go build -ldflags "-X main.GitCommit ${GIT_COMMIT}${GIT_DIRTY}")

clean:
	rm mqtt-proxy || true

.PHONY: all	
