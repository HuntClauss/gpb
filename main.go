package main

import (
	"encoding/json"
	"fmt"
	"github.com/HuntClauss/gpb/bak"
	"log"
	"os"
	"sync"
)

func main() {
	config := LoadConfig("settings.toml")
	_ = config
	client, err := bak.NewSFTPClient(config.User, config.Password, config.PhoneAddress)
	if err != nil {
		log.Fatalf("cannot create sftp client: %v\n", err)
	}
	defer client.Close()

	indexes, err := LoadIndexes(config.IndexFilePath)
	if err != nil {
		log.Fatalf("cannot load indexes: %v\n", err)
	}
	defer func() {
		if err := SaveIndexes(config.IndexFilePath, indexes); err != nil {
			fmt.Printf("cannot save indexes: %v\n", err)
		}
	}()

	var wg sync.WaitGroup
	for _, path := range config.Folders {
		log.Printf("walking folder '%s'\n", path)
		files, err := bak.Walk(client, path, indexes, config.Extensions)
		if len(files) > 0 {
			fmt.Println("Files:", len(files))
			wg.Add(1)
			go bak.CopyWorker(&wg, client, files, indexes, path, config.OutputDir)
		}

		if err != nil {
			log.Printf("folder walk for '%s' interrupted: %v\n", path, err)
		}
	}

	log.Printf("waiting for workers to finish\n")
	wg.Wait()
}

func LoadIndexes(path string) (bak.FileIndexes, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return make(bak.FileIndexes, 5000), nil
		}
		return nil, fmt.Errorf("cannot open indexes file: %w", err)
	}
	defer f.Close()

	var result bak.FileIndexes
	if err = json.NewDecoder(f).Decode(&result); err != nil {
		return nil, fmt.Errorf("cannot decode indexes file: %w", err)
	}

	return result, nil
}

func SaveIndexes(path string, indexes bak.FileIndexes) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("cannot open indexes file: %w", err)
	}

	if err = json.NewEncoder(f).Encode(indexes); err != nil {
		return fmt.Errorf("cannot encode indexes file: %w", err)
	}

	return nil
}
