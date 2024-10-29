package store

import (
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	once sync.Once
)

func ApplyMigrations(path string, dbURL string) {
	once.Do(func() {
		absPath, err := filepath.Abs(path)
		if err != nil {
			log.Fatalf("Failed to get absolute path: %v", err)
		}

		if runtime.GOOS == "windows" {
			absPath = strings.ReplaceAll(absPath, "\\", "/")
		}
		log.Printf("Using migration path: %s", absPath)
		_, err = exec.LookPath("migrate")
		if err != nil {
			log.Fatalf("migrate binary not found in PATH: %v", err)
		}

		cmd := exec.Command("migrate", "-path", absPath, "-database", dbURL, "up")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Failed to apply migrations: %v\nOutput: %s", err, output)
		}
	})
}
