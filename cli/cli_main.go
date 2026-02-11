package cli

var Main struct {
	Options struct {
		Debug bool `alt:"d" desc:"Print verbose debug information"`
	}
	Args struct{}
}

var UsesPluralFolderNames bool
