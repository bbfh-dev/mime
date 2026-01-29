package language

import (
	"fmt"

	liberrors "github.com/bbfh-dev/lib-errors"
	"github.com/bbfh-dev/mime/devkit/internal"
	"github.com/tidwall/gjson"
)

type GeneratorTemplate struct {
	Dir       string
	Iterators map[string]gjson.Result
}

func NewGeneratorTemplate(dir string, manifest *internal.JsonFile) (*GeneratorTemplate, error) {
	template := &GeneratorTemplate{Dir: dir}
	return template, nil
}

// ————————————————————————————————

type InlineTemplate struct {
	RequiredArgs []string
}

func NewInlineTemplate(dir string, manifest *internal.JsonFile) (*InlineTemplate, error) {
	template := &InlineTemplate{RequiredArgs: nil}

	if field_args := manifest.Get("arguments"); field_args.Exists() {
		switch {
		case field_args.Type == gjson.String:
			if field_args.String() != "*" {
				return nil, &liberrors.DetailedError{
					Label:   liberrors.ERR_SYNTAX,
					Context: liberrors.DirContext{Path: internal.ToAbs(dir)},
					Details: fmt.Sprintf(
						"field 'arguments' must be an array of strings or equal to '*' (string), but got %q",
						field_args.String(),
					),
				}
			}
			template.RequiredArgs = []string{}

		case field_args.IsArray():
			template.RequiredArgs = []string{}
			for _, value := range field_args.Array() {
				if value.Type != gjson.String {
					return nil, &liberrors.DetailedError{
						Label:   liberrors.ERR_SYNTAX,
						Context: liberrors.DirContext{Path: internal.ToAbs(dir)},
						Details: fmt.Sprintf(
							"field 'arguments' must be an array of strings, but got (%s) %q",
							value.Type.String(),
							value.String(),
						),
					}
				}
				template.RequiredArgs = append(template.RequiredArgs, value.String())
			}

		default:
			return nil, &liberrors.DetailedError{
				Label:   liberrors.ERR_SYNTAX,
				Context: liberrors.DirContext{Path: internal.ToAbs(dir)},
				Details: fmt.Sprintf(
					"field 'arguments' must be an object or equal to '*' (string), but got (%s) %q",
					field_args.Type.String(),
					field_args.String(),
				),
			}
		}
	}

	return template, nil
}
