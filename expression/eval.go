package expression

import (
	"fmt"
	"strings"
)

// Eval ...
func Eval(node *Node, prevCmd string, prevOutput string, siblingsOutput map[string]string, evaluator func(string) string) string {
	if node == nil {
		return ""
	}

	val := node.Source()
	if val == "" {
		return prevOutput
	}

	if prevCmd != "" {
		cmdToReplace := fmt.Sprintf("  %s ", strings.TrimSpace(prevCmd))
		val = strings.Replace(val, cmdToReplace, prevOutput, -1)
		// replace siblings output
		for cmd, output := range siblingsOutput {
			cmdToReplace := fmt.Sprintf("  %s ", strings.TrimSpace(cmd))
			val = strings.Replace(val, cmdToReplace, output, -1)
		}

	}

	sibOutput := make(map[string]string)
	if node.Siblings != nil && len(node.Siblings) > 0 {
		for _, cmd := range node.Siblings {
			sibOutput[cmd.Source()] = evaluator(cmd.Source())
		}
	}

	var nodeOutput string
	nodeOutput = evaluator(val)
	if node.Next != nil {
		nodeOutput = Eval(node.Next, node.Source(), nodeOutput, sibOutput, evaluator)
	}

	return nodeOutput
}
