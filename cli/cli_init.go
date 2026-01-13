package cli

import (
	"os"
	"path/filepath"

	libparsex "github.com/bbfh-dev/lib-parsex/v3"
	"github.com/bbfh-dev/mime/mime/errors"
	"github.com/bbfh-dev/mime/mime/minecraft"
)

var Init struct {
	Options struct {
		Name        string `alt:"n" desc:"Specify the project name that will be used for exporting" default:"untitled"`
		Minecraft   string `alt:"m" desc:"Specify the target Minecraft version. Use '-' to indicate version ranges, e.g. '1.20-1.21'" default:"1.21.11"`
		Version     string `alt:"v" desc:"Specify the project version using semantic versioning" default:"0.1.0-alpha"`
		Description string `alt:"d" desc:"Specify the project description"`
	}
	Args struct {
		WorkDir *string
	}
}

var InitProgram = libparsex.Program{
	Name:        "init",
	Description: "Initialize a new Mime project",
	Options:     &Init.Options,
	Args:        &Init.Args,
	Commands:    []*libparsex.Program{},
	EntryPoint: func(raw_args []string) error {
		if Init.Args.WorkDir != nil {
			if err := os.Chdir(*Init.Args.WorkDir); err != nil {
				return err
			}
		}

		mcmeta_body, err := os.ReadFile("pack.mcmeta")
		if err != nil {
			LogInfo(false, "Missing existing 'pack.mcmeta', so one will be created instead")
			mcmeta_body = []byte{}
		}

		mcmeta := minecraft.NewPackMcmeta(mcmeta_body)
		if value := Init.Options.Description; value != "" {
			mcmeta.File.Set("pack.description", value)
		}
		mcmeta.File.Set("meta.name", Init.Options.Name)
		mcmeta.File.Set("meta.minecraft", Init.Options.Minecraft)
		mcmeta.File.Set("meta.version", Init.Options.Version)

		err = os.WriteFile("pack.mcmeta", mcmeta.File.Formatted(), os.ModePerm)
		if err != nil {
			work_dir, _ := os.Getwd()
			return errors.NewError(
				errors.ERR_IO,
				filepath.Join(work_dir, "pack.mcmeta"),
				err.Error(),
			)
		}

		LogDone(
			false,
			"Saved 'pack.mcmeta' for name=%q version=%q minecraft=%q",
			Init.Options.Name,
			Init.Options.Version,
			Init.Options.Minecraft,
		)
		return nil
	},
}
