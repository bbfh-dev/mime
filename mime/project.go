package mime

import (
	"time"

	"github.com/bbfh-dev/mime/cli"
	"github.com/bbfh-dev/mime/mime/minecraft"
)

type Project struct {
	BuildDir           string
	Meta               *minecraft.PackMcmeta
	has_icon           bool
	has_data           bool
	has_resources      bool
	task_err           error
	data_zip_name      string
	resources_zip_name string
}

func New(mcmeta *minecraft.PackMcmeta) *Project {
	return &Project{
		BuildDir:      cli.Main.Options.Output,
		Meta:          mcmeta,
		has_icon:      false,
		has_data:      false,
		has_resources: false,
		task_err:      nil,
	}
}

func (project *Project) Build() error {
	start := time.Now()
	cli.LogInfo(
		false,
		"Building v%s for Minecraft %s",
		project.Meta.Version(),
		project.Meta.Minecraft(),
	)

	project.do(project.checkBuildDir)
	project.do(project.clearBuildDir)
	project.do(project.detectPackIcon)
	project.do(project.createResourcePack)
	project.do(project.createDataPack)

	cli.LogInfo(false, "Generating pack.mcmeta")
	if project.has_data {
		project.do(project.makePackMcmeta("data_pack", minecraft.DataPackFormats))
	}
	if project.has_resources {
		project.do(project.makePackMcmeta("resource_pack", minecraft.ResourcePackFormats))
	}

	project.do(project.runAddons)

	if project.task_err != nil {
		return project.task_err
	}
	cli.LogDone(false, "Finished building in %s", time.Since(start))

	if cli.Main.Options.Zip {
		cli.LogInfo(false, "Creating *.zip files")
		start = time.Now()

		if project.has_data {
			project.do(project.makeZip("data_pack"))
		}
		if project.has_resources {
			project.do(project.makeZip("resource_pack"))
		}

		if project.task_err != nil {
			return project.task_err
		}
		cli.LogDone(false, "Finished in %s", time.Since(start))
	}

	project.do(project.runWeld)

	return project.task_err
}

func (project *Project) do(task func() error) {
	if project.task_err == nil {
		project.task_err = task()
	}
}
