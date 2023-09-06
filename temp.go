package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	dirPrefix, _ := os.UserHomeDir()
	resultDir := filepath.Join(dirPrefix, ".L-ctl", "test")

	fmt.Printf("%s\n", resultDir)
}
