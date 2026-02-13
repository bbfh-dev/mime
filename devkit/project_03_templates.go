package devkit

func (project *Project) GenerateFromTemplates() error {
	if project.isDataCached && project.isAssetsCached {
		return nil
	}

	return nil
}
