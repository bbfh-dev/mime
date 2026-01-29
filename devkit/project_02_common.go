package devkit

import (
	"os"
	"path/filepath"

	liberrors "github.com/bbfh-dev/lib-errors"
	"github.com/bbfh-dev/mime/cli"
	"github.com/bbfh-dev/mime/devkit/internal"
	"github.com/bbfh-dev/mime/devkit/minecraft"
	cp "github.com/otiai10/copy"
	"golang.org/x/sync/errgroup"
)

func (project *Project) clearDir(path string) internal.Task {
	return func() error {
		err := os.RemoveAll(path)
		if err != nil {
			return liberrors.NewIO(err, project.BuildDir)
		}

		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return liberrors.NewIO(err, project.BuildDir)
		}

		return nil
	}
}

func (project *Project) copyPackDirs(
	folder, out_folder string,
	folders *[]string,
) internal.AsyncTask {
	return func(errs *errgroup.Group) error {
		data_entries, err := os.ReadDir(folder)
		if err != nil {
			return liberrors.NewIO(err, internal.ToAbs(folder))
		}

		for data_entry := range internal.IterateDirsOnly(data_entries) {
			path := filepath.Join(folder, data_entry.Name())
			folder_entries, err := os.ReadDir(path)
			if err != nil {
				return liberrors.NewIO(err, path)
			}

			for folder_entry := range internal.IterateDirsOnly(folder_entries) {
				path := filepath.Join(folder, data_entry.Name(), folder_entry.Name())
				switch folder_entry.Name() {
				case "function", "functions":
					if folders != nil {
						*folders = append(*folders, path)
					}
				default:
					cli.LogDebug(1, "Copying directory %q", path)
					errs.Go(func() error {
						return cp.Copy(path, filepath.Join(out_folder, path))
					})
				}
			}
		}

		return nil
	}
}

func (project *Project) copyExtraFiles(dir string) internal.Task {
	return func() error {
		for _, file := range project.extraFilesToCopy {
			cli.LogDebug(1, "Copying extra %q", file)
			path := filepath.Join(dir, file)
			err := cp.Copy(file, path)
			if err != nil {
				return liberrors.NewIO(err, path)
			}
		}
		return nil
	}
}

func (project *Project) createPackMcmeta(dir string, ft minecraft.PackFormats) internal.Task {
	return func() error {
		cli.LogInfo(1, "Exporting pack.mcmeta for %s", dir)
		mcmeta := project.Meta.Clone()
		mcmeta.FillVersion(ft)
		if err := mcmeta.SaveVersion(); err != nil {
			cli.LogWarn(1, "%s", err.Error())
		}

		path := filepath.Join(project.BuildDir, dir, "pack.mcmeta")
		err := os.WriteFile(path, mcmeta.File.Formatted(), os.ModePerm)
		if err != nil {
			return liberrors.NewIO(err, path)
		}

		return nil
	}
}
