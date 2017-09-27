package expression

import (
	"fmt"
	"strings"
)

type expander struct {
}

func (op expander) size(in []byte) int {
	if in == nil {
		return 0
	}
	count := strings.Count(string(in), "$(")
	return count
}

func (op expander) apply(in []byte) []byte {
	for i := 0; i < op.size(in); i++ {
		in = op.applyOnce(in)
	}
	return in
}

func (op expander) applyOnce(in []byte) []byte {
	if in == nil {
		return nil
	}

	counter := 0
	start := -1
	end := -1

	isValidLeft := func(index int) bool {
		a := index - 1
		b := index - 2
		if a <= 0 && b <= 0 {
			return false
		}
		return in[a] == '(' && in[b] == '$'
	}

	isValidRight := func(index int) bool {
		c := index + 1
		if c >= len(in) {
			return false
		}
		return in[c] == ')'
	}

	for i := 0; i < len(in); i++ {
		if in[i] == '`' {
			if isValidLeft(i) && counter == 0 {
				start = i
				counter++
			}
			if isValidRight(i) && counter == 1 {
				end = i
				break
			}
		}
	}

	var unroled string
	rest := in[end+1:]

	newslice := make([]byte, 0)
	if start > 0 && end != -1 && end < len(in) {
		unroled = fmt.Sprintf("filter var=%s", string(in[start+1:end]))
		newslice = append(newslice, in[0:start]...)
		newslice = append(newslice, []byte(unroled)...)
		newslice = append(newslice, rest...)
	} else {
		newslice = in
	}

	return newslice
}
