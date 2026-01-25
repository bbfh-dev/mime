package cli

var Main struct {
	Options struct {
		Output string `alt:"o" desc:"Output directory relative to the pack working dir" default:"./build"`
		Zip    bool   `alt:"z" desc:"Export data & resource packs as .zip files"`
		Debug  bool   `alt:"d" desc:"Print verbose debug information"`
		Cache  bool   `alt:"c" desc:"Use caching to prevent unnecessary computation"`
	}
	Args struct {
		WorkDir *string
	}
}
