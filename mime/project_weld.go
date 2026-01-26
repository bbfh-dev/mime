package mime

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	liberrors "github.com/bbfh-dev/lib-errors"
	"github.com/bbfh-dev/mime/cli"
	"github.com/bbfh-dev/mime/mime/internal"
	cp "github.com/otiai10/copy"
	"golang.org/x/sync/errgroup"
)

func (project *Project) weldPacks() error {
	if project.isDataCached && project.isAssetsCached {
		return nil
	}

	_, err := os.Stat("libs")
	if os.IsNotExist(err) {
		cli.LogDebug(false, "No libraries found")
		return nil
	}

	cli.LogInfo(false, "Merging with Smithed Weld")
	var errs errgroup.Group

	if !project.isDataCached {
		if _, err = os.Stat(filepath.Join("libs", "data_packs")); err == nil {
			errs.Go(func() error {
				return project.weld("data_packs", project.getZipName("DP"))
			})
		}
	}

	if !project.isAssetsCached {
		if _, err = os.Stat(filepath.Join("libs", "resource_packs")); err == nil {
			errs.Go(func() error {
				return project.weld("resource_packs", project.getZipName("RP"))
			})
		}
	}

	if err := errs.Wait(); err != nil {
		return &liberrors.DetailedError{
			Label: "Task Error",
			Context: liberrors.DirContext{
				Path: internal.ToAbs("libs"),
			},
			Details: err.Error(),
		}
	}

	return nil
}

func (project *Project) weld(dir, zip_name string) error {
	start := time.Now()
	output_name := fmt.Sprintf("weld-%s.zip", dir)

	path := filepath.Join("libs", dir)
	entries, err := readLibDir(path)
	if err != nil {
		return err
	}
	entries[len(entries)-1] = zip_name

	if len(entries) < 2 {
		cli.LogDebug(false, "No libraries found for %q", dir)
		return nil
	}

	cmd := exec.Command("weld", append([]string{
		"--dir",
		project.BuildDir,
		"--name",
		output_name,
	}, entries...)...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Stdout.Write(out.Bytes())
		return &liberrors.DetailedError{
			Label:   liberrors.ERR_EXECUTE,
			Context: liberrors.DirContext{Path: path},
			Details: err.Error(),
		}
	}

	path = filepath.Join(project.BuildDir, output_name)
	err = errors.Join(cp.Copy(path, zip_name), os.Remove(path))
	if err != nil {
		return liberrors.NewIO(err, path)
	}

	cli.LogDone(true, "Merged %q in %s", zip_name, time.Since(start))
	return nil
}

func readLibDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, liberrors.NewIO(err, internal.ToAbs("libs"))
	}

	files := make([]string, len(entries)+1)
	for i, entry := range entries {
		files[i] = filepath.Join(dir, entry.Name())
	}

	return files, nil
}
