package addons

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/bbfh-dev/mime/cli"
	"github.com/bbfh-dev/mime/mime/errors"
)

func (addon *Addon) LoadIterators() error {
	root := filepath.Join(addon.SourceDir, "iterators")
	entries, err := os.ReadDir(root)
	if err != nil {
		cli.LogDebug(true, "No iterators found")
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			cli.LogDebug(true, "Skipping directory %q", entry.Name())
			continue
		}
		name := entry.Name()

		if filepath.Ext(name) != ".csv" {
			cli.LogWarn(true, "Expected a .csv file. Skipping %q", name)
			continue
		}

		file, err := os.OpenFile(filepath.Join(root, name), os.O_RDONLY, os.ModePerm)
		if err != nil {
			return errors.NewError(errors.ERR_IO, addon.SourceDir, err.Error())
		}
		defer file.Close()

		name = trimExt(name)
		addon.Iterators[name] = Rows{}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			columns := Columns{}
			for col := range strings.SplitSeq(scanner.Text(), ",") {
				columns = append(columns, col)
			}
			addon.Iterators[name] = append(addon.Iterators[name], columns)
		}
	}

	cli.LogInfo(true, "Loaded %d iterator(s)", len(addon.Iterators))
	return nil
}

func extractIteratorsFrom(in string) []string {
	out := []string{}
	reader := bufio.NewReader(strings.NewReader(in))
	expect_bracket := false

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			return out
		}
		if char == '%' {
			expect_bracket = true
			continue
		}
		if !expect_bracket {
			continue
		}
		if expect_bracket && char != '[' {
			expect_bracket = false
			continue
		}

		identifier, err := reader.ReadString(']')
		if err != nil {
			return out
		}
		out = append(out, strings.TrimSuffix(identifier, "]"))
	}
}
