package mcfunction

import "bufio"

type McFunction struct {
	Path    string
	Scanner *bufio.Scanner
	Locals  map[string][]string
	Indent  int
}

func New(path string, scanner *bufio.Scanner, indent int) *McFunction {
	return &McFunction{
		Path:    path,
		Scanner: scanner,
		Locals:  map[string][]string{},
		Indent:  indent,
	}
}

func (fn *McFunction) Parse() error {
	// TODO: Actually parse the function
	fn.Locals[""] = []string{}
	for fn.Scanner.Scan() {
		fn.Locals[""] = append(fn.Locals[""], fn.Scanner.Text())
	}
	return Add(fn.Path, fn.Locals[""])
}
