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

	Next *Node
}

//Parse ...
func Parse(command string) *Node {
	expanded := expander{}.apply([]byte(command))
	node := parseCmd(expanded, expanded, 0)
	node = node.reverse()
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

	// Trim spaces from nodes
	for head := node; head != nil; head = head.Next {
		source := strings.TrimSpace(string(head.Src))
		head.Src = []byte(source)
		head.Start, head.End = 0, len(head.Src)-1
	}

	// Delete empty nodes: head, cursor and d
	// h = head, c = cursor, d = node to be assigned
	var h, c, d *Node = node, node.Next, nil
	for c != nil {
		if len(strings.TrimSpace(c.Source())) == 0 {
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
		fmt.Println(cursor.Source())
		cursor = cursor.Next
	}
}

// $( $( $( $() ) ) )
func parseCmd(in []byte, sub []byte, start int) *Node {
	if in == nil {
		return nil
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
			n.Next = parseCmd(in, n.Src[i:], start+i)
		}
	}
	n.End = n.Start + n.End
	n.Src = in[n.Start:n.End]

	return n
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
