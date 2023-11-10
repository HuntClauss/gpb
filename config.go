package main

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	User, Password string
	Extensions     []string
	PhoneAddress   string
	MartaMode      bool
	Folders        []string
	OutputDir      string
	IndexFilePath  string
}

func LoadConfig(path string) Config {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		log.Fatalf("cannot decode config file: %v\n", err)
	}

	return config
}
