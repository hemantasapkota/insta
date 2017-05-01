package expression

import (
	"fmt"
	"strings"
)

func Eval(node *Node, prevCmd string, prevOutput string, evaluator func(string) string) string {
	if node == nil {
		return ""
	}

	val := node.Source()

	if val == "" {
		return prevOutput
	}

	if prevCmd != "" {
		val = strings.Replace(val, fmt.Sprintf("$(%s)", prevCmd), prevOutput, -1)
	}

	var output string
	output = evaluator(val)

	if node.Next != nil {
		output = Eval(node.Next, node.Source(), output, evaluator)
	}

	return output
}
