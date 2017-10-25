package expression

import (
	"strings"
	"testing"
)

func TestCommandNode(t *testing.T) {
	var source string
	var cmd *Node
	source = `$(like id="$(last_response cmd=scrape_entry_data query=entry_data.TagPage[0].tag.media.nodes[$(counter)].id)")"`
	cmd = Parse(source)
	output := Eval(cmd, "", "", map[string]string{}, func(in string) string {
		t.Log(strings.TrimSpace(in))
		switch strings.TrimSpace(in) {
		case `counter`:
			return "0"

		case `last_response cmd=scrape_entry_data query=entry_data.TagPage[0].tag.media.nodes[0].id`:
			return "144412321231"

		case `like id="144412321231"`:
			return "Success"
		}
		return "NA"
	})
	if output != "Success" {
		t.Log("Failed eval.")
	}
}

func TestEval1(t *testing.T) {
	source := "$(filter var=result query=entry_data.TagPage[0].tag.media.nodes[$(`i`)].display_src)"
	cmd := Parse(source)
	output := Eval(cmd, "", "", map[string]string{}, func(in string) string {
		t.Log(strings.TrimSpace(in))
		switch strings.TrimSpace(in) {
		case `filter var=i`:
			return "0"

		case `filter var=result query=entry_data.TagPage[0].tag.media.nodes[0].display_src`:
			return "asdfasd"
		}
		return "NA"
	})
	if output != "asdfasd" {
		t.Error("Failed eval.")
	}
}

func TestNestedSiblings(t *testing.T) {
	source := "$(unfollow id=\"$(`followingID`)\" if=\"$(`fullName`)_contains_$(`lastName`)\")"
	cmd := Parse(source).Prune()
	if len(cmd.Siblings) != 2 {
		t.Error("Incorrect number of siblings.")
	}
	output := Eval(cmd, "", "", map[string]string{}, func(in string) string {
		source := strings.TrimSpace(in)
		t.Log(source)
		switch source {
		case `filter var=followingID`:
			return "1200"

		case `filter var=fullName`:
			return "Jon Doe"

		case `filter var=lastName`:
			return "thiki"

		case `unfollow id="1200" if="Jon Doe_contains_thiki"`:
			return "Success"

		}
		return "NA"
	})
	if output != "Success" {
		t.Error("Failed eval.")
	}
	t.Log(output)
}
