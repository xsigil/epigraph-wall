# epigraph-wall

Stealth in-memory wallpaper injector for Hyprland powered by Go, LaTeX, and `memfd_create`.

Renders classical black-letter / Fraktur epigraphs into ultra-dark desktop canvases. Zero disk I/O at runtime: embedded images are projected directly into `hyprpaper` from volatile memory.

## Architecture

1. **LaTeX (yfonts)** renders classical typographic epigraphs at 300+ DPI.
2. **ImageMagick** negates to a pitch-black dark mode canvas.
3. **Go (`go:embed`)** bakes images directly into a standalone binary.
4. **`memfd_create(2)`** creates in-memory file descriptors (`/proc/self/fd/X`) fed directly to `hyprpaper`.

## Prerequisites

- Linux (Kernel 3.17+)
- `hyprpaper` running in Hyprland
- Build-time only: `pdflatex` (with `yfonts`), `pdftoppm`, `magick`, `go`

## Build & Install

# Generates assets, builds binary, and installs to ~/.local/bin
```sh
make install
```

## Usage

# Random epigraph from RAM
```sh
epigraph-wall
```

# List embedded epigraphs
```sh
epigraph-wall -l
```

# Select specific epigraph
```sh
epigraph-wall -n 01_sovereignty.png
```

## License

MIT
