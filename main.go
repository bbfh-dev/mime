package main

import (
	"os"
	"path/filepath"

	liberrors "github.com/bbfh-dev/lib-errors"
	libparsex "github.com/bbfh-dev/lib-parsex/v3"
	"github.com/bbfh-dev/mime/cli"
	"github.com/bbfh-dev/mime/mime"
	"github.com/bbfh-dev/mime/mime/minecraft"
)

// Use -ldflags="... main.Version=<version here>"
var Version = "unset"

var MainProgram = libparsex.Program{
	Name:        "mime",
	Version:     Version,
	Description: "Minecraft data-driven vanilla data & resource pack development kit powered by pre-processors and generators",
	Options:     &cli.Main.Options,
	Args:        &cli.Main.Args,
	Commands: []*libparsex.Program{
		&cli.InitProgram,
	},
	EntryPoint: func(raw_args []string) error {
		if cli.Main.Args.WorkDir != nil {
			if err := os.Chdir(*cli.Main.Args.WorkDir); err != nil {
				return liberrors.NewIO(err, *cli.Main.Args.WorkDir)
			}
		}

		mcmeta_body, err := os.ReadFile("pack.mcmeta")
		if err != nil {
			work_dir, _ := os.Getwd()
			return liberrors.NewIO(err, work_dir)
		}

		mcmeta := minecraft.NewPackMcmeta(mcmeta_body)
		if err := mcmeta.Validate(); err != nil {
			path, _ := filepath.Abs("pack.mcmeta")
			return &liberrors.DetailedError{
				Label:   liberrors.ERR_VALIDATE,
				Context: liberrors.DirContext{Path: path},
				Details: err.Error(),
			}
		}

		project := mime.New(mcmeta)
		return project.Build()
	},
}

func main() {
	err := libparsex.Run(&MainProgram, os.Args[1:])
	if err != nil {
		switch err := err.(type) {
		case *liberrors.DetailedError:
			err.Print(os.Stderr)
		default:
			os.Stderr.WriteString(err.Error())
		}

		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
}
