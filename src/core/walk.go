package core

import (
	"os"
	"path/filepath"
)

// Walk rootPath to get all files
func Walk(rootPath string) ([]string, error) {
	files := make([]string, 0)
	if err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		files = append(files, path)
		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}
