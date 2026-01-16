package mime

import (
	"os"
	"path/filepath"

	"github.com/bbfh-dev/mime/cli"
	"github.com/bbfh-dev/mime/mime/addons"
	"github.com/bbfh-dev/mime/mime/errors"
)

func (project *Project) runAddons() error {
	_, err := os.Stat("addons")
	if os.IsNotExist(err) {
		cli.LogDebug(false, "No addons found")
		return nil
	}

	cli.LogInfo(false, "Found add-ons")

	addon_names, err := os.ReadDir("addons")
	if err != nil {
		work_dir, _ := os.Getwd()
		return errors.NewError(errors.ERR_IO, filepath.Join(work_dir, "addons"), err.Error())
	}

	for _, addon_name := range addon_names {
		if !addon_name.IsDir() {
			cli.LogDebug(true, "Skipping file %q", addon_name.Name())
			continue
		}

		cli.LogInfo(false, "Building %q", "addons/"+addon_name.Name())
		work_dir, _ := os.Getwd()
		addon := addons.New(
			filepath.Join(work_dir, "addons", addon_name.Name()),
			project.BuildDir,
		)
		if err := addon.LoadIterators(); err != nil {
			return err
		}
		if err := addon.LoadFiles("data"); err != nil {
			return err
		}
		if err := addon.LoadFiles("assets"); err != nil {
			return err
		}
		if err := addon.Build(); err != nil {
			return err
		}
	}

	return nil
}
