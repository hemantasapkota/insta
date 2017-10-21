package expression

import (
	"fmt"
	"strings"
)

//Node ...
type Node struct {
	Start int
	End   int
	Src   []byte

	Next     *Node
	Siblings []*Node
}

//Parse ...
func Parse(command string) *Node {
	expanded := expander{}.apply([]byte(command))
	siblings := make([]*Node, 0)
	node, siblings := parseCmd(expanded, expanded, 0, siblings)
	node = node.reverse()
	if len(siblings) >= 2 {
		// Remove last 2 entries ( one is empty, the other is the same as root )
		node.Siblings = siblings[0 : len(siblings)-2]
	}
	return node
}

//Source ...
func (node *Node) Source() string {
	return string(node.Src)
}

//Length ...
func (node *Node) Length() int {
	if node == nil {
		return 0
	}
	count := 0
	for cursor := node; cursor != nil; cursor = cursor.Next {
		count++
	}
	return count
}

//Prune removes all nodes that are empty
func (node *Node) Prune() *Node {
	if node == nil {
		return node
	}

	cmdHash := make(map[string]int)

	// Trim spaces from nodes
	for head := node; head != nil; head = head.Next {
		source := strings.TrimSpace(string(head.Src))
		head.Src = []byte(source)
		head.Start, head.End = 0, len(head.Src)-1
		cmdHash[source] = 1
	}

	// Count ocurrances
	for _, cmd := range node.Siblings {
		source := strings.TrimSpace(cmd.Source())
		val, ok := cmdHash[source]
		if ok {
			val = val + 1
			cmdHash[source] = val
		}
	}

	// Remove duplicates
	siblings := make([]*Node, 0)
	for i := 0; i < len(node.Siblings); i++ {
		source := strings.TrimSpace(node.Siblings[i].Source())
		if cmdHash[source] <= 1 {
			siblings = append(siblings, node.Siblings[i])
		}
	}
	node.Siblings = siblings

	// Delete empty nodes: head, cursor and d
	// h = head, c = cursor, d = node to be assigned
	var h, c, d *Node = node, node.Next, nil
	for c != nil {
		source := strings.TrimSpace(c.Source())
		if len(source) == 0 {
			d = h
			break
		}
		c = c.Next
		h = h.Next
	}
	if d != nil {
		d.Next = d.Next.Next
	}
	return node
}

//Print ...
func (node *Node) Print() {
	if node == nil {
		return
	}
	for cursor := node; cursor != nil; {
		fmt.Println(strings.TrimSpace(cursor.Source()))
		cursor = cursor.Next
	}

	fmt.Println("Siblings: ")
	for _, child := range node.Siblings {
		fmt.Println(strings.TrimSpace(child.Source()))
	}
}

// $( $( $( $() ) ) )
func parseCmd(in []byte, sub []byte, start int, children []*Node) (*Node, []*Node) {
	if in == nil {
		return nil, nil
	}

	n := &Node{Src: sub, Start: start}
	for i := 0; i < len(n.Src); i++ {
		if n.Src[i] == ')' {
			n.End = i
			n.Src[i] = ' '
			break
		}
		// Parse nested command that starts with $(
		if nextIndex := i + 1; nextIndex < len(n.Src) && n.Src[i] == '$' && n.Src[nextIndex] == '(' {
			n.Src[i], n.Src[nextIndex] = ' ', ' '
			var next *Node
			next, children = parseCmd(in, n.Src[i:], start+i, children)
			n.Next = next
		}
	}
	n.End = n.Start + n.End
	n.Src = in[n.Start:n.End]
	children = append(children, n)

	return n, children
}

func (node *Node) reverse() *Node {
	doReverse := func() *Node {
		var head, next, cursor *Node = node, nil, nil
		for head != nil {
			next = head.Next
			head.Next = cursor
			cursor = head
			head = next
		}
		return cursor
	}
	return doReverse()
}
