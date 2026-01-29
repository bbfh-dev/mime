package devkit

import (
	"github.com/bbfh-dev/mime/cli"
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

	for name, template := range project.generatorTemplates {
		cli.LogDebug(1, "Generating from %q", name)

		for name := range template.Definitions {
			cli.LogDebug(2, "Create %q", name)
		}

		cli.LogDone(1, "Finished generating %q for %d definitions", name, len(template.Definitions))
	}

	return nil
}
