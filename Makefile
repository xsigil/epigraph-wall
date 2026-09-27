PREFIX ?= $(HOME)/.local
BINDIR = $(PREFIX)/bin
MANDIR = $(PREFIX)/share/man/man1

BIN = epigraph-wall
MAN = doc/epigraph-wall.1
TEX_SRCS = $(wildcard tex/*.tex)

.PHONY: all assets build install uninstall clean

all: build

assets: $(TEX_SRCS)
	@echo "==> Compiling LuaLaTeX epigraphs into images/..."
	./scripts/build-assets.sh

build: assets
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
