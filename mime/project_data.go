package mime

import (
	"os"
	"path/filepath"

	liberrors "github.com/bbfh-dev/lib-errors"
	"github.com/bbfh-dev/mime/cli"
	"github.com/bbfh-dev/mime/mime/internal"
	"golang.org/x/sync/errgroup"

	cp "github.com/otiai10/copy"
)

func (project *Project) genDataPack() error {
	if project.isDataCached {
		return nil
	}

	_, err := os.Stat("data")
	if os.IsNotExist(err) {
		cli.LogDebug(false, "No data pack found")
		return nil
	}

	cli.LogInfo(false, "Creating a data pack")
	path := filepath.Join(project.BuildDir, "data_pack")

	err = os.RemoveAll(path)
	if err != nil {
		return liberrors.NewIO(err, project.BuildDir)
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return liberrors.NewIO(err, project.BuildDir)
	}

	data_entries, err := os.ReadDir("data")
	if err != nil {
		return liberrors.NewIO(err, internal.ToAbs("data"))
	}

	funcFoldersToParse := []string{}
	var errs errgroup.Group

	for data_entry := range internal.IterateDirsOnly(data_entries) {
		path = filepath.Join("data", data_entry.Name())
		folder_entries, err := os.ReadDir(path)
		if err != nil {
			return liberrors.NewIO(err, path)
		}

		for folder_entry := range internal.IterateDirsOnly(folder_entries) {
			path = filepath.Join("data", data_entry.Name(), folder_entry.Name())
			switch folder_entry.Name() {
			case "function", "functions":
				funcFoldersToParse = append(funcFoldersToParse, path)
			default:
				cli.LogDebug(true, "Copying directory %q", path)
				errs.Go(func() error {
					return cp.Copy(path, filepath.Join(project.BuildDir, "data_pack", path))
				})
			}
		}
	}

	// for _, path = range funcFoldersToParse {
	// 	filepath.WalkDir(path, func(path string, entry fs.DirEntry, err error) error {
	// 		if err != nil || entry.IsDir() {
	// 			return err
	// 		}
	// 		errs.Go(func() error {
	// 			// return project.parseFunction(path)
	// 			return nil
	// 		})
	// 		return nil
	// 	})
	// }

	if err := errs.Wait(); err != nil {
		return &liberrors.DetailedError{
			Label: "Task Error",
			Context: liberrors.DirContext{
				Path: internal.ToAbs("data"),
			},
			Details: err.Error(),
		}
	}

	for _, file := range project.extraFilesToCopy {
		cli.LogDebug(true, "Copying extra %q", file)
		path = filepath.Join(project.BuildDir, "data_pack", file)
		err = cp.Copy(file, path)
		if err != nil {
			return liberrors.NewIO(err, path)
		}
	}

	return nil
}
