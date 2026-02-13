package code_test

import (
	"testing"

	"github.com/bbfh-dev/vintage/devkit/internal/code"
	"github.com/tidwall/gjson"
	"gotest.tools/assert"
)

func TestSubstituteString(t *testing.T) {
	env := code.NewEnv()
	env.Variables["test"] = gjson.Parse(`{"nested": {"within": 123}}`)
	result, err := code.SubstituteString("Hello %[test.nested.within]!", env)
	assert.NilError(t, err)
	assert.DeepEqual(t, result, "Hello 123!")

	env.Variables["test2"] = code.SimpleVariable("World")
	result, err = code.SubstituteString("Hello %[test2]!", env)
	assert.NilError(t, err)
	assert.DeepEqual(t, result, "Hello World!")
}
