package mime

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bbfh-dev/mime/cli"
	"github.com/bbfh-dev/mime/mime/errors"
	"github.com/bbfh-dev/mime/mime/minecraft"
)

func (project *Project) checkBuildDir() error {
	project.BuildDir, _ = filepath.Abs(project.BuildDir)
	cli.LogInfo(false, "Checking build directory: %s", project.BuildDir)

	stat, err := os.Stat(project.BuildDir)
	if err != nil {
		if os.IsNotExist(err) {
			cli.LogDebug(true, "Directory doesn't exist yet. Skipping checks")
			return nil
		}
		return errors.NewError(errors.ERR_IO, project.BuildDir, err.Error())
	}

	if !stat.IsDir() {
		return errors.NewError(
			errors.ERR_VALID,
			project.BuildDir,
			"build output is a file",
		)
	}

	return nil
}

func (project *Project) clearBuildDir() error {
	cli.LogInfo(false, "Clearing build directory")

	if err := os.RemoveAll(project.BuildDir); err != nil {
		return errors.NewError(errors.ERR_IO, project.BuildDir, err.Error())
	}

	err := os.MkdirAll(project.BuildDir, os.ModePerm)
	if err != nil {
		return errors.NewError(errors.ERR_IO, project.BuildDir, err.Error())
	}

	return nil
}

func (project *Project) detectPackIcon() error {
	if project.isCached() {
		return nil
	}

	_, err := os.Stat("pack.png")
	if os.IsNotExist(err) {
		cli.LogWarn(false, "No pack icon found")
		return nil
	}

	cli.LogInfo(false, "Found 'pack.png'")
	project.has_icon = true
	return nil
}

func (project *Project) makePackMcmeta(
	name string,
	formats map[string]minecraft.Version,
) func() error {
	return func() error {
		cli.LogInfo(true, "Exporting into the %s", name)
		mcmeta := minecraft.NewPackMcmeta(project.Meta.File.Body)
		mcmeta.FillVersion(formats)
		if err := mcmeta.SaveVersion(); err != nil {
			cli.LogWarn(true, "%s", err.Error())
		}

		path := filepath.Join(project.BuildDir, name, "pack.mcmeta")
		err := os.WriteFile(path, mcmeta.File.Formatted(), os.ModePerm)
		if err != nil {
			return errors.NewError(errors.ERR_IO, path, err.Error())
		}

		return nil
	}
}

func (project *Project) getZipName(label string) string {
	return filepath.Join(
		project.BuildDir,
		fmt.Sprintf(
			"%s_%s_v%s.zip",
			project.Meta.Name(),
			label,
			project.Meta.PrintableVersion(),
		),
	)
}

func (project *Project) makeZip(folder, cache_name string) func() error {
	return func() error {
		if folder != "data_pack" && folder != "resource_pack" {
			panic("Folder must only be data_pack or resource_pack")
		}
		path := project.getZipName(strings.ToUpper(string(folder[0])) + "P")

		if slices.Contains(project.cached, cache_name) {
			return nil
		}

		file, err := os.Create(path)
		if err != nil {
			return errors.NewError(errors.ERR_IO, path, err.Error())
		}
		defer file.Close()

		writer := zip.NewWriter(file)
		defer writer.Close()

		path = filepath.Join(project.BuildDir, folder)
		root, err := os.OpenRoot(path)
		if err != nil {
			return errors.NewError(errors.ERR_IO, path, err.Error())
		}

		err = writer.AddFS(root.FS())
		if err != nil {
			return errors.NewError(errors.ERR_IO, path, err.Error())
		}

		return nil
	}
}
