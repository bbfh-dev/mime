package devkit

import (
	"os"
	"path/filepath"

	liberrors "github.com/bbfh-dev/lib-errors"
	liblog "github.com/bbfh-dev/lib-log"
	cp "github.com/otiai10/copy"
)

func (project *Project) CopyOverlays() error {
	liblog.Debug(0, "Copying overlays...")

	entries, err := os.ReadDir("overlays")
	if err != nil {
		return liberrors.NewIO(err, "overlays")
	}

	for _, entry := range entries {
		for _, group := range [][]string{{"data_pack", "data"}, {"resource_pack", "assets"}} {
			target, folder := group[0], group[1]
			src := filepath.Join("overlays", entry.Name(), folder)
			if exists(src) {
				err := cp.Copy(
					src,
					filepath.Join(project.BuildDir, target, "overlays", entry.Name(), folder),
				)
				if err != nil {
					return liberrors.NewIO(err, src)
				}
				liblog.Done(1, "Copied %s", entry.Name())
			} else {
				err := os.MkdirAll(
					filepath.Join(project.BuildDir, target, "overlays", entry.Name()),
					os.ModePerm,
				)
				if err != nil {
					return liberrors.NewIO(err, src)
				}
			}
		}
	}

	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
