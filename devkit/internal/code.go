package internal

import "io"

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
