package language

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/bbfh-dev/mime/devkit/internal"
)

type Mcfunction struct {
	Path    string
	Scanner *bufio.Scanner

	root    *Line
	current *Line
}

func NewMcfunction(path string, scanner *bufio.Scanner) *Mcfunction {
	root := newLine(nil, "<root>")
	return &Mcfunction{
		Path:    path,
		Scanner: scanner,
		root:    root,
		current: root,
	}
}

func (fn *Mcfunction) BuildTree() *Mcfunction {
	indents := []int{0}

	for fn.Scanner.Scan() {
		parent_indent := indents[len(indents)-1]

		line := fn.Scanner.Text()
		line_indent := internal.GetIndentOf(line) - parent_indent
		line = strings.TrimSpace(line)

		if line_indent > 0 {
			indents = append(indents, parent_indent+line_indent)
			fn.current = fn.current.Nested[len(fn.current.Nested)-1]
			goto next_iteration
		}

		if line_indent < 0 {
			indents = indents[:len(indents)-1]
			fn.current = fn.current.Parent
			goto next_iteration
		}

	next_iteration:
		fn.current.Append(newLine(fn.current, line))
	}

	return fn
}

func (mcfunction *Mcfunction) Parse() error {
	fmt.Println(mcfunction.root)
	return nil
}

// ————————————————————————————————

type Line struct {
	Parent   *Line
	Contents string
	Nested   []*Line
}

func newLine(parent *Line, contents string) *Line {
	return &Line{Parent: parent, Contents: contents, Nested: nil}
}

func (line *Line) Append(nested *Line) {
	line.Nested = append(line.Nested, nested)
}

func (line *Line) Write(writer io.Writer, indent int) {
	writer.Write([]byte(line.Contents + "\n"))
	for _, line := range line.Nested {
		internal.WriteIndentString(writer, indent+4)
		line.Write(writer, indent+4)
	}
}

func (line *Line) String() string {
	var builder strings.Builder
	line.Write(&builder, 0)
	return builder.String()
}
