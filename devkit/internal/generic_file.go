package internal

type GenericFile struct {
	Ext  string
	Body []byte
}

func NewGenericFile(ext string, body []byte) *GenericFile {
	return &GenericFile{
		Ext:  ext,
		Body: body,
	}
}

func (file *GenericFile) Clone() *GenericFile {
	return NewGenericFile(file.Ext, file.Body)
}

func (file *GenericFile) Formatted() []byte {
	switch file.Ext {

	case ".json":
		return NewJsonFile(file.Body).Formatted()

	default:
		return file.Body
	}
}
