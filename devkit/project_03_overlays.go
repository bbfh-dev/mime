package devkit

import (
	"os"
	"path/filepath"

	"codeberg.org/bbfh/vintage/devkit/internal/pipeline"
	liberrors "github.com/bbfh-dev/lib-errors"
	liblog "github.com/bbfh-dev/lib-log"
	cp "github.com/otiai10/copy"
	"golang.org/x/sync/errgroup"
)

func (project *Project) LoadOverlays() error {
	path := "overlays"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return liberrors.NewIO(err, path)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dir := filepath.Join(path, entry.Name(), "data")
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			project.dataOverlays = append(project.dataOverlays, entry.Name())
		}

		dir = filepath.Join(path, "assets")
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			project.assetsOverlays = append(project.assetsOverlays, entry.Name())
		}
	}

	liblog.Info(0, "Found %d overlay(s)", len(project.dataOverlays)+len(project.assetsOverlays))
	return nil
}

func (project *Project) copyOverlays(source []string, target, folder string) pipeline.AsyncTask {
	if len(source) == 0 {
		return nil
	}

	return func(errs *errgroup.Group) error {
		for _, overlay := range source {
			path := filepath.Join("overlays", overlay, folder)

			liblog.Debug(0, "Copying overlay %s", overlay)
			errs.Go(func() error {
				return cp.Copy(
					path,
					filepath.Join(project.BuildDir, target, overlay, folder),
				)
			})
		}

		return nil
	}
}
