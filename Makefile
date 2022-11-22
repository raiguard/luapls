all: luapls

luapls:
	@go build -o luapls

clean:
	rm -f luapls

test: *.go
	@go test

.PHONY: clean luapls test
