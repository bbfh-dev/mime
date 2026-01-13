package mime

import (
	"os"
	"path/filepath"

	"github.com/bbfh-dev/mime/cli"
	"github.com/bbfh-dev/mime/mime/errors"
	cp "github.com/otiai10/copy"
)

func (project *Project) createResourcePack() error {
	_, err := os.Stat("assets")
	if os.IsNotExist(err) {
		cli.LogDebug(false, "No resource pack found")
		return nil
	}

	cli.LogInfo(false, "Creating a resource pack")
	project.has_resources = true

	err = os.MkdirAll(filepath.Join(project.BuildDir, "resource_pack"), os.ModePerm)
	if err != nil {
		return errors.NewError(errors.ERR_IO, project.BuildDir, err.Error())
	}

	cli.LogInfo(true, "Copying 'assets/'")
	path := filepath.Join(project.BuildDir, "resource_pack", "assets")
	err = cp.Copy("assets", path)
	if err != nil {
		return errors.NewError(errors.ERR_IO, path, err.Error())
	}

	if project.has_icon {
		path := filepath.Join(project.BuildDir, "resource_pack", "pack.png")
		err = cp.Copy("pack.png", path)
		if err != nil {
			return errors.NewError(errors.ERR_IO, path, err.Error())
		}
	}

	return nil
}
