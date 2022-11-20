luapls: *.go lua/*.go
	@go build -o luapls

clean:
	rm -f luapls

test: *.go
	@go test

.PHONY: clean
