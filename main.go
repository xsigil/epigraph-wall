package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

//go:embed images/*.png
var embeddedImages embed.FS

func main() {
	listFlag := flag.Bool("l", false, "List embedded epigraph images")
	nameFlag := flag.String("n", "", "Specify epigraph image name")
	monitorFlag := flag.String("m", "", "Specify monitor name (leave empty for all)")
	flag.Parse()

	entries, err := fs.Glob(embeddedImages, "images/*.png")
	if err != nil || len(entries) == 0 {
		fmt.Fprintln(os.Stderr, "[-] Error: No embedded images found in binary.")
		os.Exit(1)
	}

	if *listFlag {
		fmt.Println("Embedded Epigraphs:")
		for _, e := range entries {
			fmt.Printf("  - %s\n", filepath.Base(e))
		}
		return
	}

	selected := entries[0]
	if *nameFlag != "" {
		found := false
		for _, e := range entries {
			if filepath.Base(e) == *nameFlag {
				selected = e
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "[-] Epigraph '%s' not found.\n", *nameFlag)
			os.Exit(1)
		}
	} else {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		selected = entries[r.Intn(len(entries))]
	}

	data, err := embeddedImages.ReadFile(selected)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[-] Read embedded error: %v\n", err)
		os.Exit(1)
	}

	// XDG_RUNTIME_DIR (RAM 上の tmpfs) に一時展開
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = "/tmp"
	}

	tmpFile := filepath.Join(runtimeDir, fmt.Sprintf("epigraph-%d.png", os.Getpid()))
	if err := os.WriteFile(tmpFile, data, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "[-] Failed to write to tmpfs: %v\n", err)
		os.Exit(1)
	}
	// 適用後に RAM から即時消去
	defer os.Remove(tmpFile)

	fmt.Printf("[*] Injecting %s via RAM tmpfs (%s)\n", filepath.Base(selected), tmpFile)

	// 新構文: hyprctl hyprpaper wallpaper "<monitor>,<path>"
	wallArg := fmt.Sprintf("%s,%s", *monitorFlag, tmpFile)
	wall := exec.Command("hyprctl", "hyprpaper", "wallpaper", wallArg)
	if out, err := wall.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "[-] hyprpaper wallpaper failed: %s (%v)\n", string(out), err)
		os.Exit(1)
	}

	// hyprpaper の GPU テクスチャロード完了を待機
	time.Sleep(500 * time.Millisecond)
	fmt.Println("[+] Wallpaper applied seamlessly from RAM.")
}
