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
