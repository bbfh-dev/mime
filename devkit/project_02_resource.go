package devkit

import (
	"os"
	"path/filepath"

	"github.com/bbfh-dev/mime/cli"
	"github.com/bbfh-dev/mime/devkit/internal"
	"github.com/bbfh-dev/mime/devkit/minecraft"
)

func (project *Project) GenerateResourcePack() error {
	if project.isAssetsCached {
		return nil
	}

	_, err := os.Stat(FOLDER_ASSETS)
	if os.IsNotExist(err) {
		cli.LogDebug(0, "No resource pack found")
		return nil
	}

	cli.LogInfo(0, "Creating a Resource Pack")
	path := filepath.Join(project.BuildDir, "resource_pack")

	return internal.Pipeline(
		project.clearDir(path),
		internal.Async(
			project.copyPackDirs(FOLDER_ASSETS, path, nil),
		),
		project.copyExtraFiles(path),
		project.createPackMcmeta("resource_pack", minecraft.ResourcePackFormats),
	)
}
