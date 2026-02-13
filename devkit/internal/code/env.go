package code

import "github.com/tidwall/gjson"

type Columns []string
type Rows []Columns

type Env struct {
	Iterators map[string]Columns
	Variables map[string]Variable
}

func NewEnv() Env {
	return Env{
		Iterators: map[string]Columns{},
		Variables: map[string]Variable{},
	}
}

type Variable interface {
	String() string
	Value() any
}

func IsStringifiable(variable Variable) bool {
	if variable, ok := variable.(gjson.Result); ok {
		return variable.Type != gjson.JSON
	}
	return true
}

func TypeOf(variable Variable) gjson.Type {
	if variable, ok := variable.(gjson.Result); ok {
		return variable.Type
	}
	return gjson.String
}

func Query(variable Variable, path string) Variable {
	if variable, ok := variable.(gjson.Result); ok {
		return variable.Get(path)
	}
	return gjson.Result{Type: gjson.Null}
}

type SimpleVariable string

func (variable SimpleVariable) String() string {
	return string(variable)
}

func (variable SimpleVariable) Value() any {
	return string(variable)
}
