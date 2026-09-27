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

	"golang.org/x/sys/unix"
)

//go:embed images/*.png
var embeddedImages embed.FS

func createMemFile(name string, data []byte) (*os.File, error) {
	fd, err := unix.MemfdCreate(name, unix.MFD_CLOEXEC)
	if err != nil {
		return nil, fmt.Errorf("memfd_create failed: %w", err)
	}

	file := os.NewFile(uintptr(fd), name)
	if _, err := file.Write(data); err != nil {
		file.Close()
		return nil, fmt.Errorf("write to memfd failed: %w", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		file.Close()
		return nil, fmt.Errorf("seek memfd failed: %w", err)
	}

	return file, nil
}

func main() {
	listFlag := flag.Bool("l", false, "List embedded epigraph images")
	nameFlag := flag.String("n", "", "Specify epigraph image name")
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

	memFile, err := createMemFile("epigraph.png", data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[-] Memory mapping failed: %v\n", err)
		os.Exit(1)
	}
	defer memFile.Close()

	memPath := fmt.Sprintf("/proc/%d/fd/%d", os.Getpid(), memFile.Fd())
	fmt.Printf("[*] Injecting %s via RAM (%s)\n", filepath.Base(selected), memPath)

	preload := exec.Command("hyprctl", "hyprpaper", "preload", memPath)
	if out, err := preload.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "[-] hyprpaper preload failed: %s (%v)\n", string(out), err)
		os.Exit(1)
	}

	wall := exec.Command("hyprctl", "hyprpaper", "wallpaper", fmt.Sprintf(",%s", memPath))
	if out, err := wall.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "[-] hyprpaper wallpaper failed: %s (%v)\n", string(out), err)
		os.Exit(1)
	}

	// hyprpaper の非同期ロード完了を待機
	time.Sleep(1 * time.Second)
	fmt.Println("[+] Wallpaper applied seamlessly from RAM.")
}
