package expression

import (
	"testing"
)

func TestCommandNode(t *testing.T) {

	var source string = ""
	var cmd *Node

	// Test 1
	source = `$(repeat frequency=10 cmd="$(like id=$(last_response cmd=scrape_entry_data query=entry_data.TagPage[0].tag.media.nodes[$(counter)].id)))"`
	cmd = Parse(source)
	// cmd.Print()

	// Test 2
	source = `$(like id="$(last_response cmd=scrape_entry_data query=entry_data.TagPage[0].tag.media.nodes[$(counter)].id)")"`
	cmd = Parse(source)
	cmd.Print()

	output := Eval(cmd, "", "", func(in string) string {
		switch in {
		case `counter`:
			return "0"

		case `last_response cmd=scrape_entry_data query=entry_data.TagPage[0].tag.media.nodes[0].id`:
			return "144412321231"

		case `like id="144412321231"`:
			return "success"
		}
		return "NA"
	})

	t.Log("Output is ", output)

	source = `$(last_response cmd=scrape_entry_data query=entry_data.TagPage[0].tag.media.nodes[$(counter)].id)`
	cmd = Parse(source)
	cmd.Print()
}
