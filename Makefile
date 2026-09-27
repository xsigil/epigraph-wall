PREFIX ?= $(HOME)/.local
BINDIR = $(PREFIX)/bin
MANDIR = $(PREFIX)/share/man/man1

BIN = epigraph-wall
MAN = doc/epigraph-wall.1

.PHONY: all assets build install uninstall clean

all: build

assets:
	@echo "==> Rendering LaTeX epigraphs into images/..."
	./scripts/gen-assets.sh

build:
	@if [ -z "$$(ls -A images/*.png 2>/dev/null)" ]; then \
		$(MAKE) assets; \
	fi
	@echo "==> Building static binary with embedded images..."
	go build -ldflags="-s -w" -o $(BIN) main.go

install: build
	@mkdir -p $(BINDIR) $(MANDIR)
	install -m 755 $(BIN) $(BINDIR)/$(BIN)
	install -m 644 $(MAN) $(MANDIR)/epigraph-wall.1

uninstall:
	rm -f $(BINDIR)/$(BIN)
	rm -f $(MANDIR)/epigraph-wall.1

clean:
	rm -f $(BIN)
	rm -f images/*.png
