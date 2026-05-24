package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	BufoDir    string
	StoreFile  string
	SocketFile string
	PIDFile    string
	ProxyPort  int
}

func Load() *Config {
	bufoDir := getBufoDir()

	return &Config{
		BufoDir:    bufoDir,
		StoreFile:  filepath.Join(bufoDir, StoreFile),
		SocketFile: filepath.Join(bufoDir, SocketFile),
		PIDFile:    filepath.Join(bufoDir, PIDFile),
		ProxyPort:  ProxyPort,
	}
}

func getBufoDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(homeDir, ".bufo")
}
