package main

import (
	"os"

	libparsex "github.com/bbfh-dev/lib-parsex/v3"
)

const VERSION = "0.1.0-alpha.1"

var Main struct {
	Options struct {
		Output string `alt:"o" desc:"Output directory relative to the pack working dir" default:"./build"`
		Zip    bool   `alt:"z" desc:"Export data & resource packs as .zip files"`
		Force  bool   `alt:"f" desc:"Force building even if the build directory looks off"`
		Debug  bool   `alt:"d" desc:"Print verbose debug information"`
	}
	Args struct {
		WorkDir *string
	}
}

var MainProgram = libparsex.Program{
	Name:        "mime",
	Version:     VERSION,
	Description: "Minecraft data & resource pack processor designed to not significantly modify the Minecraft syntax while providing useful code generation",
	Options:     &Main.Options,
	Args:        &Main.Args,
	Commands:    []*libparsex.Program{},
	EntryPoint: func(raw_args []string) error {
		if Main.Args.WorkDir != nil {
			if err := os.Chdir(*Main.Args.WorkDir); err != nil {
				return err
			}
		}

		return nil
	},
}

func main() {
	err := libparsex.Run(&MainProgram, os.Args[1:])
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
