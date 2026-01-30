package devkit

import (
	"fmt"
	"path/filepath"

	"github.com/bbfh-dev/mime/cli"
	"github.com/bbfh-dev/mime/devkit/internal"
)

func (project *Project) GenerateFromTemplates() error {
	if project.isDataCached && project.isAssetsCached {
		return nil
	}

	if len(project.generatorTemplates) == 0 {
		cli.LogDebug(0, "No generator templates defined")
		return nil
	}

	cli.LogInfo(0, "Generating code from %d template(s)", len(project.generatorTemplates))

	for template_name, template := range project.generatorTemplates {
		cli.LogDebug(1, "Generating from %q", template_name)

		for definition_name, definition := range template.Definitions {
			cli.LogDebug(2, "Create %q", template_name)

			root := filepath.Join("templates", template_name)
			tree, err := internal.LoadTree(
				root,
				[2]string{"data", "data_pack"},
				[2]string{"assets", "resource_pack"},
			)
			if err != nil {
				return err
			}

			for path, file := range tree {
				fmt.Println(path, file, definition_name, definition)
			}
		}

		cli.LogDone(
			1,
			"Finished generating %q for %d definitions",
			template_name,
			len(template.Definitions),
		)
	}

	return nil
}
