package mime

import (
	"time"

	"github.com/bbfh-dev/mime/cli"
	"github.com/bbfh-dev/mime/mime/minecraft"
)

type Project struct {
	Meta *minecraft.PackMcmeta
}

func New(mcmeta *minecraft.PackMcmeta) *Project {
	return &Project{
		Meta: mcmeta,
	}
}

func (project *Project) Build() error {
	start := time.Now()
	cli.LogInfo(
		false,
		"Building v%s for Minecraft %s",
		project.Meta.File.Get("meta.version"),
		project.Meta.File.Get("meta.minecraft"),
	)

	cli.LogDone(false, "Finished building in %s", time.Since(start))
	return nil
}
