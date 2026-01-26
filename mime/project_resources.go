package mime

import (
	"os"

	"github.com/bbfh-dev/mime/cli"
)

func (project *Project) genResourcePack() error {
	if project.isAssetsCached {
		return nil
	}

	_, err := os.Stat("assets")
	if os.IsNotExist(err) {
		cli.LogDebug(false, "No resource pack found")
		return nil
	}

	return nil
}
