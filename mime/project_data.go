package mime

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bbfh-dev/mime/cli"
	"github.com/bbfh-dev/mime/mime/errors"
	"github.com/bbfh-dev/mime/mime/mcfunction"
	cp "github.com/otiai10/copy"
	"golang.org/x/sync/errgroup"
)

func (project *Project) createDataPack() error {
	if project.isDataCached() {
		cli.LogDebug(true, "The data pack is cached")
		return nil
	}

	_, err := os.Stat("data")
	if os.IsNotExist(err) {
		cli.LogDebug(false, "No data pack found")
		return nil
	}

	cli.LogInfo(false, "Creating a data pack")
	project.has_data = true

	err = os.MkdirAll(filepath.Join(project.BuildDir, "data_pack"), os.ModePerm)
	if err != nil {
		return errors.NewError(errors.ERR_IO, project.BuildDir, err.Error())
	}

	data_entries, err := os.ReadDir("data")
	if err != nil {
		work_dir, _ := os.Getwd()
		return errors.NewError(errors.ERR_IO, filepath.Join(work_dir, "data"), err.Error())
	}

	function_paths := []string{}
	var errs errgroup.Group

	for _, data_entry := range data_entries {
		if !data_entry.IsDir() {
			cli.LogDebug(true, "Skipping file %q", data_entry.Name())
			continue
		}

		folder_entries, err := os.ReadDir(filepath.Join("data", data_entry.Name()))
		if err != nil {
			work_dir, _ := os.Getwd()
			return errors.NewError(
				errors.ERR_IO,
				filepath.Join(work_dir, "data", data_entry.Name()),
				err.Error(),
			)
		}

		for _, folder_entry := range folder_entries {
			path := filepath.Join("data", data_entry.Name(), folder_entry.Name())
			if !folder_entry.IsDir() {
				cli.LogDebug(true, "Skipping file %q", path)
				continue
			}

			switch folder_entry.Name() {
			case "function", "functions":
				function_paths = append(function_paths, path)
			default:
				cli.LogDebug(true, "Copying directory %q", path)
				errs.Go(func() error {
					return cp.Copy(
						path,
						filepath.Join(
							project.BuildDir,
							"data_pack",
							"data",
							data_entry.Name(),
							folder_entry.Name(),
						),
					)
				})
			}
		}
	}

	if err := errs.Wait(); err != nil {
		work_dir, _ := os.Getwd()
		return errors.NewError(
			errors.ERR_INTERNAL,
			filepath.Join(work_dir, "data"),
			err.Error(),
		)
	}

	cli.LogInfo(true, "Parsing mcfunction files")
	for _, path := range function_paths {
		filepath.WalkDir(path, func(path string, entry fs.DirEntry, err error) error {
			if err != nil || entry.IsDir() {
				return err
			}

			errs.Go(func() error {
				return project.parseFunction(path)
			})

			return nil
		})
	}

	if err := errs.Wait(); err != nil {
		work_dir, _ := os.Getwd()
		return errors.NewError(
			errors.ERR_INTERNAL,
			filepath.Join(work_dir, "data"),
			err.Error(),
		)
	}

	cli.LogInfo(true, "Writing mcfunction files to disk")

	for path, lines := range mcfunction.Registry {
		path = filepath.Join(project.BuildDir, "data_pack", path)

		err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			return errors.NewError(errors.ERR_IO, path, err.Error())
		}

		err = os.WriteFile(path, []byte(strings.Join(lines, "\n")), os.ModePerm)
		if err != nil {
			return errors.NewError(errors.ERR_IO, path, err.Error())
		}
	}

	if project.has_icon {
		path := filepath.Join(project.BuildDir, "data_pack", "pack.png")
		err = cp.Copy("pack.png", path)
		if err != nil {
			return errors.NewError(errors.ERR_IO, path, err.Error())
		}
	}

	return nil
}

func (project *Project) parseFunction(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return errors.NewError(errors.ERR_IO, path, err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	function := mcfunction.New(path, scanner, 0)
	return function.Parse()
}
