package devkit

import (
	"os"
	"path/filepath"

	"codeberg.org/bbfh/vintage/devkit/internal/pipeline"
	"codeberg.org/bbfh/vintage/devkit/minecraft"
	liblog "github.com/bbfh-dev/lib-log"
)

func (project *Project) GenerateResourcePack() error {
	if project.isAssetsCached {
		return nil
	}

	_, err := os.Stat(FOLDER_ASSETS)
	if os.IsNotExist(err) {
		liblog.Debug(0, "No resource pack found")
		return nil
	}

	liblog.Info(0, "Creating a Resource Pack")
	path := filepath.Join(project.BuildDir, "resource_pack")

	return pipeline.New(
		project.clearDir(path),
		pipeline.Async(
			project.copyPackDirs(FOLDER_ASSETS, path, nil),
		),
		pipeline.Async(
			pipeline.If[pipeline.AsyncTask](project.hasPatches()).
				Then(project.patch("resource_pack", "assets")),
		),
		project.copyExtraFiles(path),
		pipeline.Async(
			project.copyOverlays(project.assetsOverlays, "resource_pack", "assets"),
		),
		project.createPackMcmeta(
			"resource_pack",
			"resources",
			minecraft.ResourcePackFormats,
			project.assetsOverlays,
		),
	)
}
