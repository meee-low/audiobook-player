package main

import (
	"fmt"
	"io"
	"log"
	"os"

	toml "github.com/pelletier/go-toml/v2"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <config_file_path>", os.Args[0])

	}

	fp := os.Args[1]
	log.Printf("File path: %s", fp)

	_, err := os.Stat(fp)

	if os.IsNotExist(err) {
		writeDefaultConfig(fp)
	}

	file, err := os.Open(fp)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	config := NewConfig()

	config.load_config(file)

	config.print_config()
}

type Config struct {
	Directories          []string `toml:"directories"`
	DatabasePath         string   `toml:"database_path"`
	Keybinds             keybinds `toml:"keybinds"`
	DefaultPlaybackSpeed float32  `toml:"default_playback_speed"`
	AutoResume           bool     `toml:"auto_resume"`
}

type keybinds struct {
	Pause       string `toml:"pause"`
	Play        string `toml:"play"`
	PlayOrPause string `toml:"play_or_pause"`
	Background  string `toml:"background"`
}

func NewConfig() Config {
	return Config{
		Directories:          []string{},
		DatabasePath:         "",
		Keybinds:             newKeybinds(),
		DefaultPlaybackSpeed: 1.0,
		AutoResume:           false,
	}
}

func newKeybinds() keybinds {
	return keybinds{
		Pause:       "",
		Play:        "",
		PlayOrPause: "",
		Background:  "",
	}
}

func (c *Config) load_config(config_file *os.File) {
	data, err := io.ReadAll(config_file)
	if err != nil {
		log.Fatalf("Failed to read config file: %s\n", err)
	}
	// read toml file:
	err = toml.Unmarshal(data, c)
	if err != nil {
		log.Fatalf("Error while unmarshaling the file: %s\n", err)
	}
}

func (c *Config) print_config() {
	fmt.Printf("Directories: %v\n", c.Directories)
	fmt.Printf("Database path: %s\n", c.DatabasePath)
	fmt.Printf("Autoresume: %t\n", c.AutoResume)
	fmt.Printf("Default playback speed: %.2f\n", c.DefaultPlaybackSpeed)
	c.Keybinds.print_keybinds()
}

func (kb *keybinds) print_keybinds() {
	fmt.Printf("pause: %s\n", kb.Pause)
	fmt.Printf("play: %s\n", kb.Play)
	fmt.Printf("playOrPause: %s\n", kb.PlayOrPause)
	fmt.Printf("background: %s\n", kb.Background)
}

func writeDefaultConfig(config_path string) {
	config := NewConfig()
	data, err := toml.Marshal(config)
	if err != nil {
		log.Fatalf("Failed to marshal the new config: %s", err)
	}

	err = os.WriteFile(config_path, data, 0o644)
	if err != nil {
		log.Fatalf("Failed to write default config file: %s", err)
	} else {
		log.Printf("Default config file created at %s", config_path)
	}
}
