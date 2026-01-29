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

	// for name, template := range project.generatorTemplates {
	// }

	return nil
}
