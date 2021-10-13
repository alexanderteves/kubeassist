.PHONY: all test

all: test

test:
	go test -count=1 ./...