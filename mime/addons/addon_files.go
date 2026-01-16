package addons

import (
	"io/fs"
	"os"
	"path/filepath"
)

func (addon *Addon) LoadFiles(name string) error {
	path := filepath.Join(addon.SourceDir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	return filepath.WalkDir(path, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}

		addon.Paths = append(addon.Paths, path)
		return nil
	})
}
