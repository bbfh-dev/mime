package patcher

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	liberrors "github.com/bbfh-dev/lib-errors"
)

const MAX_SEEK_OFFSET = 200

var hunkHeaderRe = regexp.MustCompile(`^@@\s+-(\d+)(?:,(\d+))?\s+\+(\d+)(?:,(\d+))?\s+@@`)

type Line struct {
	Op   byte   // ' ', '-', '+'
	Text string // without the leading op character
}

type Hunk struct {
	OldStart, OldCount int
	NewStart, NewCount int
	Lines              []Line
}

func Patch(patch, path string) error {
	hunks, err := parseHunks(patch)
	if err != nil {
		return &liberrors.DetailedError{
			Label:   liberrors.ERR_FORMAT,
			Details: err.Error(),
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return liberrors.NewIO(err, path)
	}

	// Preserve whether the original file ended with a newline
	endsWithNL := bytes.HasSuffix(data, []byte{'\n'})

	scanner := bufio.NewScanner(bytes.NewReader(data))
	var fileLines []string
	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return &liberrors.DetailedError{
			Label:   liberrors.ERR_READ,
			Details: err.Error(),
		}
	}

	newLines, err := applyHunks(fileLines, hunks)
	if err != nil {
		return &liberrors.DetailedError{
			Label:   liberrors.ERR_WRITE,
			Details: err.Error(),
		}
	}

	var buf bytes.Buffer
	for i, line := range newLines {
		buf.WriteString(line)
		if i < len(newLines)-1 || endsWithNL {
			buf.WriteByte('\n')
		}
	}

	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return liberrors.NewIO(err, path)
	}

	return nil
}

func parseHunks(patch string) ([]Hunk, error) {
	var hunks []Hunk
	scanner := bufio.NewScanner(strings.NewReader(patch))

	var current *Hunk

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "@@") {
			m := hunkHeaderRe.FindStringSubmatch(line)
			if m == nil {
				return nil, fmt.Errorf("invalid hunk header: %q", line)
			}

			oldStart, _ := strconv.Atoi(m[1])
			oldCount := 1
			if m[2] != "" {
				oldCount, _ = strconv.Atoi(m[2])
			}
			newStart, _ := strconv.Atoi(m[3])
			newCount := 1
			if m[4] != "" {
				newCount, _ = strconv.Atoi(m[4])
			}

			h := Hunk{
				OldStart: oldStart,
				OldCount: oldCount,
				NewStart: newStart,
				NewCount: newCount,
			}
			hunks = append(hunks, h)
			current = &hunks[len(hunks)-1]
			continue
		}

		if current == nil {
			// skip any leading noise
			continue
		}

		if len(line) == 0 {
			// empty line is treated as context
			current.Lines = append(current.Lines, Line{Op: ' ', Text: ""})
			continue
		}

		op := line[0]
		text := line[1:]

		switch op {
		case ' ', '-', '+':
			current.Lines = append(current.Lines, Line{Op: op, Text: text})
		case '\\':
			// "\ No newline at end of file" – ignore for now
		default:
			return nil, fmt.Errorf("unexpected line in hunk: %q", line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(hunks) == 0 {
		return nil, fmt.Errorf("no hunks found in patch")
	}
	return hunks, nil
}

// findHunkLocation tries the recorded position first, then searches
// nearby for a matching context (offset recovery).
func findHunkLocation(fileLines []string, h Hunk) (int, error) {
	// Convert 1-based to 0-based
	want := h.OldStart - 1

	// Build the expected old lines (context + deleted)
	var expected []string
	for _, l := range h.Lines {
		if l.Op == ' ' || l.Op == '-' {
			expected = append(expected, l.Text)
		}
	}

	if len(expected) == 0 {
		// pure insertion – just trust the line number
		if want < 0 {
			want = 0
		}
		if want > len(fileLines) {
			want = len(fileLines)
		}
		return want, nil
	}

	matchesAt := func(pos int) bool {
		if pos < 0 || pos+len(expected) > len(fileLines) {
			return false
		}
		for i, exp := range expected {
			if fileLines[pos+i] != exp {
				return false
			}
		}
		return true
	}

	// 1. Try exact location
	if matchesAt(want) {
		return want, nil
	}

	// 2. Search with increasing offset
	for offset := 1; offset <= MAX_SEEK_OFFSET; offset++ {
		// forward
		if matchesAt(want + offset) {
			return want + offset, nil
		}
		// backward
		if matchesAt(want - offset) {
			return want - offset, nil
		}
	}

	return -1, fmt.Errorf(
		"context mismatch for hunk @@ -%d,%d +%d,%d @@",
		h.OldStart,
		h.OldCount,
		h.NewStart,
		h.NewCount,
	)
}

func applyHunks(fileLines []string, hunks []Hunk) ([]string, error) {
	// We apply hunks from last to first so earlier line numbers stay valid
	// while we are still searching.
	for _, h := range slices.Backward(hunks) {
		pos, err := findHunkLocation(fileLines, h)
		if err != nil {
			return nil, err
		}

		// Build the replacement block (context + added lines)
		var replacement []string
		for _, l := range h.Lines {
			if l.Op == ' ' || l.Op == '+' {
				replacement = append(replacement, l.Text)
			}
		}

		// How many old lines does this hunk consume?
		oldLen := 0
		for _, l := range h.Lines {
			if l.Op == ' ' || l.Op == '-' {
				oldLen++
			}
		}

		// Splice
		newLines := make([]string, 0, len(fileLines)-oldLen+len(replacement))
		newLines = append(newLines, fileLines[:pos]...)
		newLines = append(newLines, replacement...)
		newLines = append(newLines, fileLines[pos+oldLen:]...)
		fileLines = newLines
	}

	return fileLines, nil
}
