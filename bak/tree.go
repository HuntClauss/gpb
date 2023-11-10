package bak

import (
	"fmt"
	"github.com/pkg/sftp"
	"os"
	"strings"
)

type FileIndexes map[string]int64

type FileInfo struct {
	AbsolutePath string
	FileStat     os.FileInfo
}

func Walk(client *sftp.Client, path string, fileIndexes FileIndexes, exts []string) ([]FileInfo, error) {
	result := make([]FileInfo, 0, 5000)

	counter := 0
	walker := client.Walk(path)
	for walker.Step() {
		if !hasWhitlistedExt(walker.Stat().Name(), exts) {
			continue
		}

		if val, ok := fileIndexes[walker.Path()]; ok && val == walker.Stat().ModTime().Unix() {
			continue
		}

		walker.Stat().Mode().Type()
		result = append(result, FileInfo{
			AbsolutePath: walker.Path(),
			FileStat:     walker.Stat(),
		})
		counter++
		//if counter >= 20 {
		//	break
		//}
	}

	if walker.Err() != nil {
		return nil, fmt.Errorf("cannot walk: %v\n", walker.Err())
	}

	return result, nil
}

func hasWhitlistedExt(path string, exts []string) bool {
	for _, ext := range exts {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}
