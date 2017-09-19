package expression

import (
	"strings"
	"testing"
)

func TestCommandNode(t *testing.T) {
	var source string
	var cmd *Node

	// Test 2
	source = `$(like id="$(last_response cmd=scrape_entry_data query=entry_data.TagPage[0].tag.media.nodes[$(counter)].id)")"`
	cmd = Parse(source)

	output := Eval(cmd, "", "", func(in string) string {
		t.Log(strings.TrimSpace(in))
		switch strings.TrimSpace(in) {
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
}

func TestEval1(t *testing.T) {
	source := "$(filter var=result query=entry_data.TagPage[0].tag.media.nodes[$(`i`)].display_src)"
	cmd := Parse(source)

	output := Eval(cmd, "", "", func(in string) string {
		t.Log(strings.TrimSpace(in))
		switch strings.TrimSpace(in) {
		case `filter var=i`:
			return "0"

		case `filter var=result query=entry_data.TagPage[0].tag.media.nodes[0].display_src`:
			return "asdfasd"
		}
		return "NA"
	})

	t.Log(output)
}
