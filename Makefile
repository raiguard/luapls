PREFIX = $(DESTDIR)/usr/local
BINDIR = $(PREFIX)/bin

all: luapls

luapls:
	@go build -o luapls

clean:
	rm -f luapls

test: *.go
	@go test ./...

install:
	install -Dpm 0755 luapls $(BINDIR)/luapls

uninstall:
	rm -f $(BINDIR)/luapls

.PHONY: clean luapls test install uninstall
