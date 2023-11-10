package bak

import (
	"github.com/pkg/sftp"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func CopyWorker(wg *sync.WaitGroup, client *sftp.Client, files []FileInfo, indexes FileIndexes, rootDir, outputDir string) {
	log.Printf("copying files from '%s' to '%s' started\n", rootDir, outputDir)
	for _, file := range files {
		rel, err := filepath.Rel("/", file.AbsolutePath)
		if err != nil {
			log.Printf("cannot get relative path for '%s': %v\n", file.AbsolutePath, err)
			continue
		}

		folder := filepath.Dir(rel)
		filename := filepath.Base(rel)

		targetFolder := filepath.Join(outputDir, folder)
		if err = os.MkdirAll(targetFolder, 0744); err != nil {
			log.Printf("cannot create folder '%s': %v\n", targetFolder, err)
			continue
		}

		sourceFile, err := client.OpenFile(file.AbsolutePath, os.O_RDONLY)
		if err != nil {
			log.Printf("cannot open remote file '%s': %v\n", file.AbsolutePath, err)
			continue
		}

		targetFilePath := filepath.Join(outputDir, folder, filename)
		targetFile, err := os.OpenFile(targetFilePath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			_ = sourceFile.Close()
			log.Printf("cannot open local file '%s': %v\n", targetFilePath, err)
			continue
		}

		if _, err = io.Copy(targetFile, sourceFile); err != nil {
			_, _ = targetFile.Close(), sourceFile.Close()
			log.Printf("cannot copy file '%s' to '%s': %v\n", file.AbsolutePath, targetFilePath, err)
			continue
		}
		log.Printf("copied '%s' to '%s'\n", file.AbsolutePath, targetFilePath)

		indexes[file.AbsolutePath] = file.FileStat.ModTime().Unix()
		_, _ = targetFile.Close(), sourceFile.Close()
		if err = os.Chtimes(targetFilePath, file.FileStat.ModTime(), file.FileStat.ModTime()); err != nil {
			log.Printf("cannot change '%s' mtime and atime: %v\n", targetFolder, err)
			continue
		}
	}

	log.Printf("copying files from '%s' to '%s' finished\n", rootDir, outputDir)
	wg.Done()
}
