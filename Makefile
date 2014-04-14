GOPATH := `echo $$GOPATH`

all: clean
	$(shell cd content; $(GOPATH)/bin/rice embed)
	$(shell go build)

clean:
	rm mqtt-proxy || true

.PHONY: all	
