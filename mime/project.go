package mime

import (
	"os"
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
	cached             []string
}

func New(mcmeta *minecraft.PackMcmeta) *Project {
	return &Project{
		BuildDir:      cli.Main.Options.Output,
		Meta:          mcmeta,
		has_icon:      false,
		has_data:      false,
		has_resources: false,
		task_err:      nil,
		cached:        []string{},
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

	project.data_zip_name = project.getZipName("DP")
	project.resources_zip_name = project.getZipName("RP")

	project.do(project.checkBuildDir)
	project.do(project.clearBuildDir)
	project.do(project.detectPackIcon)
	if cli.Main.Options.Cache {
		var addons_time time.Time
		if _, err := os.Stat("addons"); err == nil {
			addons_time = getLastModified("addons")
		}
		project.do(project.initCache("data", addons_time))
		project.do(project.initCache("assets", addons_time))
	}
	project.do(project.createResourcePack)
	project.do(project.createDataPack)

	cli.LogInfo(false, "Generating pack.mcmeta")
	if project.has_data && !project.isDataCached() {
		project.do(project.makePackMcmeta("data_pack", minecraft.DataPackFormats))
	}
	if project.has_resources && !project.isResourcesCached() {
		project.do(project.makePackMcmeta("resource_pack", minecraft.ResourcePackFormats))
	}

	project.do(project.runAddons)

	if project.task_err != nil {
		return project.task_err
	}
	if !project.isCached() {
		cli.LogDone(false, "Finished building in %s", time.Since(start))
	}

	if cli.Main.Options.Zip {
		cli.LogInfo(false, "Creating *.zip files")
		start = time.Now()

		if project.has_data && !project.isDataCached() {
			project.do(project.makeZip("data_pack", "data"))
		}
		if project.has_resources && !project.isResourcesCached() {
			project.do(project.makeZip("resource_pack", "assets"))
		}

		if project.task_err != nil {
			return project.task_err
		}
		if !project.isCached() {
			cli.LogDone(false, "Finished in %s", time.Since(start))
		}
	}

	project.do(project.runWeld)
	if cli.Main.Options.Cache {
		project.do(project.saveCache("data", project.data_zip_name))
		project.do(project.saveCache("assets", project.resources_zip_name))
	}

	return project.task_err
}

func (project *Project) do(task func() error) {
	if project.task_err == nil {
		project.task_err = task()
	}
}
