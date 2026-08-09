package devkit

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"codeberg.org/bbfh/vintage/cli"
	"codeberg.org/bbfh/vintage/devkit/internal/patcher"
	"codeberg.org/bbfh/vintage/devkit/internal/pipeline"
	liberrors "github.com/bbfh-dev/lib-errors"
	liblog "github.com/bbfh-dev/lib-log"
	cp "github.com/otiai10/copy"
	"golang.org/x/sync/errgroup"
)

func (project *Project) patch(target, folder string) pipeline.AsyncTask {
	return func(errs *errgroup.Group) error {
		names := strings.Split(cli.Build.Options.Patches, "+")
		liblog.Info(1, "Found %d patch candidates", len(names))

		for _, name := range names {
			if _, err := os.Stat(filepath.Join("patches", name)); os.IsNotExist(err) {
				liblog.Warn(1, "Unknown overlay %q. Skipping...", name)
				continue
			}

			root := filepath.Join("patches", name, folder)
			if _, err := os.Stat(root); err != nil {
				continue
			}

			liblog.Debug(2, "Applying %s patch", name)
			err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return err
				}

				relPath := path
				relPath = strings.TrimPrefix(path, root)
				relPath = filepath.Join(folder, relPath)
				destPath := filepath.Join(project.BuildDir, target, relPath)

				switch filepath.Ext(path) {

				case ".patch":
					patch, err := os.ReadFile(path)
					if err != nil {
						return liberrors.NewIO(err, path)
					}
					destPath = strings.TrimSuffix(destPath, ".patch")
					if err := patcher.Patch(
						string(patch),
						destPath,
					); err != nil {
						return err
					}
					liblog.Done(2, "Patched %s", destPath)

				default:
					liblog.Debug(2, "Copying %s", path)
					if err := cp.Copy(path, destPath); err != nil {
						return liberrors.NewIO(err, path)
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func (project *Project) hasPatches() bool {
	_, err := os.Stat("patches")
	return err == nil
}

// TODO: this
// versions := strings.SplitN(entry.Name(), "-", 2)
// if len(versions) != 2 {
// 	return &liberrors.DetailedError{
// 		Label:   liberrors.ERR_FORMAT,
// 		Context: nil,
// 		Details: "Overlay name must be [mc_version]__[mc_version], but got: " + entry.Name(),
// 	}
// }
//
// min_version_id := versions[0]
// max_version_id := versions[1]
// for _, version := range []string{min_version_id, max_version_id} {
// 	if !minecraft.IsVersionSupported(version) {
// 		return &liberrors.DetailedError{
// 			Label:   liberrors.ERR_VALIDATE,
// 			Context: nil,
// 			Details: "Overlay minecraft version is not supported, or most likely not written in the correct format: " + version,
// 		}
// 	}
// }
// min_version := versionProvider[min_version_id]
// max_version := versionProvider[max_version_id]
