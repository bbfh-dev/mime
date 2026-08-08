package devkit

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"codeberg.org/bbfh/vintage/cli"
	"codeberg.org/bbfh/vintage/devkit/internal/pipeline"
	liberrors "github.com/bbfh-dev/lib-errors"
	liblog "github.com/bbfh-dev/lib-log"
	cp "github.com/otiai10/copy"
)

func (project *Project) ApplyPatches(folder, target string) pipeline.Task {
	return func() error {
		liblog.Debug(0, "Applying patches...")

		for patch := range strings.SplitSeq(cli.Build.Options.Patches, "+") {
			liblog.Info(0, "Applying patch %q", patch)

			path := filepath.Join("patches", patch)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				liblog.Warn(0, "Patch %q is not found. Skipping...", patch)
				continue
			}

			dir := filepath.Join(path, folder)
			_, err := os.Stat(dir)
			if err != nil {
				continue
			}

			err = project.applyPatch(dir, folder, target)
			if err != nil {
				return err
			}
		}

		liblog.Done(0, "Patched %s", target)
		return nil
	}
}

func (project *Project) applyPatch(path, folder, target string) error {
	liblog.Debug(1, "Applying patches from %s", path)

	filepath.WalkDir(path, func(local string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		dest := filepath.Join(project.BuildDir, target, folder, strings.TrimPrefix(local, path))
		switch filepath.Ext(local) {
		case ".patch":
			body, err := os.ReadFile(local)
			if err != nil {
				return liberrors.NewIO(err, local)
			}

			dest = strings.TrimSuffix(dest, filepath.Ext(dest))
			contents := strings.ReplaceAll(string(body), "${VINTAGE_TARGET}", dest)

			cmd := exec.Command("patch", "-l", "-p0", "-d", "/", "--no-backup-if-mismatch")
			cmd.Stdin = strings.NewReader(contents)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return &liberrors.DetailedError{
					Label:   liberrors.ERR_EXECUTE,
					Context: liberrors.NewProgramContext(cmd, ""),
					Details: err.Error(),
				}
			}

		default:
			liblog.Debug(1, "Copying %s", local)
			if err := cp.Copy(local, dest); err != nil {
				return liberrors.NewIO(err, local)
			}
		}
		return nil
	})

	return nil
}
