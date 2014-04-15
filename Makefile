GOPATH := `echo $$GOPATH`

# Get the git commit
GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_DIRTY := $(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)

all: clean
	$(shell cd content; $(GOPATH)/bin/rice embed)
	$(shell go build -ldflags "-X main.GitCommit ${GIT_COMMIT}${GIT_DIRTY}")

clean:
	rm mqtt-proxy || true

.PHONY: all	
