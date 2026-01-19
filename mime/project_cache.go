package mime

import (
	"encoding/binary"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/bbfh-dev/mime/cli"
	"github.com/bbfh-dev/mime/mime/errors"
	"github.com/emersion/go-appdir"
	cp "github.com/otiai10/copy"
)

func (project *Project) isDataCached() bool {
	return slices.Contains(project.cached, "data")
}

func (project *Project) isResourcesCached() bool {
	return slices.Contains(project.cached, "assets")
}

func (project *Project) isCached() bool {
	return len(project.cached) > 0
}

func (project *Project) initCache(folder string, addons_time time.Time) func() error {
	return func() error {
		latest_time := getLastModified(folder)
		if addons_time.Sub(latest_time) > 0 {
			latest_time = addons_time
		}

		cache_dir := appdir.New("mime").UserCache()
		if err := os.MkdirAll(cache_dir, os.ModePerm); err != nil {
			work_dir, _ := os.Getwd()
			return errors.NewError(errors.ERR_IO, filepath.Join(work_dir, folder), err.Error())
		}

		path := filepath.Join(cache_dir, fmt.Sprintf("%s_timestamp.bin", folder))
		data, err := os.ReadFile(path)
		if err != nil {
			cli.LogDebug(true, "%q is not cached", folder)
			return nil
		}

		timestamp := time.UnixMilli(int64(binary.BigEndian.Uint64(data)))

		if latest_time.Sub(timestamp) < 0 {
			cli.LogInfo(true, "%q is cached", folder)
			project.cached = append(project.cached, folder)
			return nil
		}

		cli.LogDebug(true, "%q cache is outdated", folder)
		return nil
	}
}

func (project *Project) saveCache(folder, zip_name string) func() error {
	return func() error {
		cache_dir := appdir.New("mime").UserCache()
		if slices.Contains(project.cached, folder) {
			path := filepath.Join(cache_dir, folder+".zip")
			if err := cp.Copy(path, zip_name); err != nil {
				return errors.NewError(errors.ERR_IO, path, err.Error())
			}
			cli.LogDone(false, "Loaded cached %q", zip_name)
			return nil
		}

		stat, err := os.Stat(zip_name)
		if err != nil {
			cli.LogDebug(true, "ERROR: %s", err.Error())
			return nil
		}

		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(stat.ModTime().UnixMilli()))

		err = os.WriteFile(
			filepath.Join(cache_dir, fmt.Sprintf("%s_timestamp.bin", folder)),
			buffer,
			os.ModePerm,
		)
		if err != nil {
			work_dir, _ := os.Getwd()
			return errors.NewError(errors.ERR_IO, filepath.Join(work_dir, folder), err.Error())
		}

		err = cp.Copy(zip_name, filepath.Join(cache_dir, folder+".zip"))
		if err != nil {
			work_dir, _ := os.Getwd()
			return errors.NewError(errors.ERR_IO, filepath.Join(work_dir, folder), err.Error())
		}

		cli.LogDebug(true, "Cached %q", zip_name)
		return nil
	}
}

func getLastModified(dir string) time.Time {
	var latest_time time.Time
	filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil && entry.IsDir() {
			return err
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		if timestamp := info.ModTime(); timestamp.Sub(latest_time) > 0 {
			latest_time = timestamp
		}

		return nil
	})
	return latest_time
}
