package templates

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	liberrors "github.com/bbfh-dev/lib-errors"
	"github.com/bbfh-dev/vintage/devkit/internal/code"
	"github.com/bbfh-dev/vintage/devkit/internal/drive"
	"github.com/tidwall/gjson"
)

const BODY_SUBSTITUTION = "%[...]"
const SNIPPET_FILENAME = "snippet.mcfunction"

type InlineTemplate struct {
	RequiredArgs []string
	Call         func(out io.Writer, in io.Reader, args []string) error
}

func NewInlineTemplate(dir string, manifest *drive.JsonFile) (*InlineTemplate, error) {
	template := &InlineTemplate{RequiredArgs: nil}

	field_args := manifest.Get("arguments")
	if field_args.Exists() {
		switch {

		case field_args.IsArray():
			template.RequiredArgs = []string{}
			for _, value := range field_args.Array() {
				if value.Type != gjson.String {
					return nil, newSyntaxError(
						drive.ToAbs(dir),
						"field 'arguments' must be an array of strings",
						value,
					)
				}
				template.RequiredArgs = append(template.RequiredArgs, value.String())
			}

		default:
			return nil, newSyntaxError(
				drive.ToAbs(dir),
				"field 'arguments' must be an array of strings",
				field_args,
			)
		}
	}

	path := filepath.Join(dir, SNIPPET_FILENAME)
	if _, err := os.Stat(path); err == nil {
		return inlineTemplateUsingSnippet(template, path)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, liberrors.NewIO(err, dir)
	}

	for entry := range drive.IterateFilesOnly(entries) {
		switch {
		case strings.HasPrefix(entry.Name(), "call"):
			path := filepath.Join(dir, entry.Name())
			return inlineTemplateUsingExec(template, path)
		}
	}

	return template, &liberrors.DetailedError{
		Label:   liberrors.ERR_VALIDATE,
		Context: liberrors.DirContext{Path: drive.ToAbs(dir)},
		Details: fmt.Sprintf(
			"template %q contains no logic files. Must contain `*.mcfunction` or `call*`. Refer to documentation",
			filepath.Base(dir),
		),
	}
}

func (template *InlineTemplate) IsArgPassthrough() bool {
	return template.RequiredArgs == nil
}

func newSyntaxError(path, details string, field gjson.Result) *liberrors.DetailedError {
	return &liberrors.DetailedError{
		Label:   liberrors.ERR_SYNTAX,
		Context: liberrors.DirContext{Path: path},
		Details: fmt.Sprintf("%s, but got (%s) %q", details, field.Type, field),
	}
}

func inlineTemplateUsingSnippet(template *InlineTemplate, path string) (*InlineTemplate, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, liberrors.NewIO(err, drive.ToAbs(path))
	}

	template.Call = func(writer io.Writer, reader io.Reader, args []string) error {
		env := code.NewEnv()
		for i, arg := range args {
			env.Variables[template.RequiredArgs[i]] = code.SimpleVariable(arg)
		}

		body := string(body)
		lines := strings.Split(body, "\n")

		var before strings.Builder
		var after string
		var ok bool
		for i, line := range lines {
			if strings.Contains(line, BODY_SUBSTITUTION) {
				after = strings.Join(lines[i+1:], "\n")
				ok = true
				break
			}
			before.WriteString(line + "\n")
		}

		if ok {
			if err := writeSubstituted(writer, path, before.String(), env); err != nil {
				return err
			}
			io.Copy(writer, reader)
			if err := writeSubstituted(writer, path, after, env); err != nil {
				return err
			}
		}

		return nil
	}

	return template, nil
}

func inlineTemplateUsingExec(template *InlineTemplate, path string) (*InlineTemplate, error) {
	template.Call = func(writer io.Writer, reader io.Reader, args []string) error {
		cmd := exec.Command(path, args...)
		cmd.Stdin = reader
		cmd.Stdout = writer
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			return &liberrors.DetailedError{
				Label:   liberrors.ERR_EXECUTE,
				Context: liberrors.DirContext{Path: path},
				Details: err.Error(),
			}
		}

		return nil
	}

	return template, nil
}

func writeSubstituted(writer io.Writer, path, in string, env code.Env) error {
	str, err := code.SubstituteString(in, env)
	if err != nil {
		return &liberrors.DetailedError{
			Label:   liberrors.ERR_FORMAT,
			Context: liberrors.DirContext{Path: path},
			Details: err.Error(),
		}
	}
	writer.Write([]byte(str))
	return nil
}
