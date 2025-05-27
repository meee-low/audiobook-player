package config

// TODO: Read the config file from the env variables, or set to a default path (based on the OS).

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

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

func (c *Config) LoadConfig(config_file *os.File) {
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

func (c *Config) PrintConfig() {
	fmt.Printf("Directories: %v\n", c.Directories)
	fmt.Printf("Database path: %s\n", c.DatabasePath)
	fmt.Printf("Autoresume: %t\n", c.AutoResume)
	fmt.Printf("Default playback speed: %.2f\n", c.DefaultPlaybackSpeed)
	c.Keybinds.PrintKeybinds()
}

func (kb *keybinds) PrintKeybinds() {
	fmt.Printf("pause: %s\n", kb.Pause)
	fmt.Printf("play: %s\n", kb.Play)
	fmt.Printf("playOrPause: %s\n", kb.PlayOrPause)
	fmt.Printf("background: %s\n", kb.Background)
}

func WriteDefaultConfig(config_path string) {
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

func DefaultConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Could not determine home directory: %s", err)
	}
	return filepath.Join(homeDir, ".config", "audiobook-player")
}

func DefaultConfigPath() string {
	return filepath.Join(DefaultConfigDir(), "config.toml")
}

func LoadOrCreateConfig(config_path string) Config {
	// Check if the config file exists
	_, err := os.Stat(config_path)
	if os.IsNotExist(err) {
		log.Printf("Config file does not exist at %s, creating default config.", config_path)
		WriteDefaultConfig(config_path)
	}

	file, err := os.Open(config_path)
	if err != nil {
		log.Fatalf("Error opening config file: %s", err)
	}
	defer file.Close()

	config := NewConfig()
	config.LoadConfig(file)

	return config
}
