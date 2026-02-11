package cli

import (
	"os"

	liberrors "github.com/bbfh-dev/lib-errors"
)

var Main struct {
	Options struct {
		Debug bool `alt:"d" desc:"Print verbose debug information"`
	}
	Args struct{}
}

var UsesPluralFolderNames bool

func ApplyWorkDir(work_dir *string) error {
	if work_dir != nil {
		if err := os.Chdir(*work_dir); err != nil {
			return liberrors.NewIO(err, *work_dir)
		}
	}
	return nil
}
