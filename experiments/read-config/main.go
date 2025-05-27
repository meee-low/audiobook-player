package main

import (
	"log"
	"os"

	"github.com/meee-low/audiobook-player/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <config_file_path>", os.Args[0])
	}

	fp := os.Args[1]
	log.Printf("File path: %s", fp)

	config := config.LoadOrCreateConfig(fp)

	config.PrintConfig()
}
