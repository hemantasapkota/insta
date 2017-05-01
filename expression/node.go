package expression

import (
	"fmt"
)

// Note: Command should start with $()
// Parses: $(repeat frequency=10 cmd="$(like id=$(last_response cmd=scrape_entry_data query=entry_data.TagPage[0].tag.media.nodes[$(counter)].id)))

// Into a stack of:
// repeat frequency=10 cmd="$(like id=$(last_response cmd=scrape_entry_data query=entry_data.TagPage[0].tag.media.nodes[$(counter)].id))
// like id=$(last_response cmd=scrape_entry_data query=entry_data.TagPage[0].tag.media.nodes[$(counter)].id)
// last_response cmd=scrape_entry_data query=entry_data.TagPage[0].tag.media.nodes[$(counter)].id
// counter

type Node struct {
	Start int
	End   int
	Src   []byte

	Next *Node
}

func Parse(command string) *Node {
	node := parseCmd([]byte(command), []byte(command), 0)
	node = node.reverse()
	return node
}

func (node *Node) Source() string {
	return string(node.Src)
}

func (node *Node) Print() {
	if node == nil {
		return
	}

	for cursor := node; cursor != nil; {
		fmt.Println(string(cursor.Src))
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

		if nextIndex := i + 1; nextIndex < len(n.Src) && n.Src[i] == '$' && n.Src[nextIndex] == '(' {
			n.Src[i] = ' '
			n.Src[i+1] = ' '
			from := i + 2

			n.Next = parseCmd(in, n.Src[from:], start+from)
		}
	}

	n.End = n.Start + n.End
	n.Src = in[n.Start:n.End]

	return n
}

func (node *Node) reverse() *Node {
	doReverse := func() *Node {
		// Init our three pointers
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
