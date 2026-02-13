package code

import (
	"bufio"
	"strings"
)

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
