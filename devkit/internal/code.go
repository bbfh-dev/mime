package internal

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"
)

func GetIndentOf(line string) (indent int) {
	for _, char := range line {
		switch char {
		case ' ':
			indent++
		case '\t':
			indent += 4
		default:
			return
		}
	}
	return
}

func WriteIndentString(writer io.Writer, indent int) {
	if indent%4 == 0 {
		char := []byte{'\t'}
		for range indent / 4 {
			writer.Write(char)
		}
		return
	}

	char := []byte{' '}
	for range indent {
		writer.Write(char)
	}
}

func GetIndentString(indent int) string {
	var builder strings.Builder
	WriteIndentString(&builder, indent)
	return builder.String()
}

func ExtractVariablesFrom(in string) []string {
	out := []string{}
	reader := bufio.NewReader(strings.NewReader(in))
	expect_bracket := false

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			return out
		}

		switch {

		case char == '%':
			expect_bracket = true

		case !expect_bracket:
			// ignore

		case expect_bracket && char != '[':
			expect_bracket = false

		default:
			identifier, err := reader.ReadString(']')
			if err != nil {
				return out
			}
			out = append(out, strings.TrimSuffix(identifier, "]"))
		}
	}
}

func PathToResource(path string) string {
	fields := strings.Split(path, "/")

	// Convert into pack local space
	if index := slices.Index(fields, "data"); index != -1 {
		fields = fields[index+1:]
	} else if index := slices.Index(fields, "assets"); index != -1 {
		fields = fields[index+1:]
	}

	switch len(fields) {
	case 0, 1, 2:
		panic(fmt.Sprintf(
			"Invalid PathToResource(%q). Not enough directories to convert",
			path,
		))
	default:
		last := len(fields) - 1
		fields[last] = strings.TrimSuffix(
			fields[last],
			filepath.Ext(fields[last]),
		)
		return fields[0] + ":" + strings.Join(fields[2:], "/")
	}
}

func ResourceToPath(folder_name, resource string) string {
	parts := strings.SplitN(resource, ":", 2)
	if len(parts) == 1 {
		panic(fmt.Sprintf("Invalid ResourceToFilepath(%q, %q)", folder_name, resource))
	}
	return filepath.Join(parts[0], folder_name, parts[1])
}

func Fields(in string) []string {
	reader := bufio.NewReader(strings.NewReader(in))
	var out []string
	var builder strings.Builder

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			goto exit
		}

		switch char {

		case ' ':
			out = append(out, builder.String())
			builder.Reset()

		case '"', '\'', '`':
			str, err := reader.ReadString(byte(char))
			if err != nil {
				goto exit
			}
			if char == '`' {
				builder.WriteString(strings.TrimSuffix(str, string(char)))
			} else {
				builder.WriteString(string(char) + str)
			}

		default:
			builder.WriteRune(char)
		}
	}

exit:
	if builder.Len() != 0 {
		out = append(out, builder.String())
	}
	return out
}
