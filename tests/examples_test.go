package vintage_test

import (
	"os"
	"path/filepath"
	"testing"

	"codeberg.org/bbfh/vintage/cli"
	"codeberg.org/bbfh/vintage/devkit"
	liblog "github.com/bbfh-dev/lib-log"
	"gotest.tools/assert"
)

func TestExamples(t *testing.T) {
	work_dir, err := os.Getwd()
	assert.NilError(t, err)

	work_dir = filepath.Join(work_dir, "..")
	assert.NilError(t, os.Chdir(work_dir))

	entries, err := os.ReadDir("examples")
	assert.NilError(t, err)

	cli.Main.Options.Debug = testing.Verbose()
	liblog.LogLevel = liblog.LEVEL_DEBUG
	cli.Build.Options.Force = true
	cli.Build.Options.Output = filepath.Join(os.TempDir(), "vintage-test")
	cli.Build.Options.Zip = true

	for _, entry := range entries {
		path := filepath.Join("examples", entry.Name())

		// Reset
		devkit.Reset()
		assert.NilError(t, os.RemoveAll(cli.Build.Options.Output))
		assert.NilError(t, os.Chdir(work_dir))

		t.Run(entry.Name(), func(t *testing.T) {
			liblog.Output = t.Output()
			cli.Build.Args.WorkDir = &path
			if entry.Name() == "05_patches" {
				cli.Build.Options.Patches = "example1+example2"
			}
			err := devkit.Build([]string{path})
			assert.NilError(t, err)
		})
	}
}
