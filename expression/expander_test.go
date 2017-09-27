package expression

import (
	"testing"
)

func TestExpander(t *testing.T) {

	// download url="$(`url`)" file="0.jpg"
	in := []byte("download url=\"$(`url`)\"")
	out := []byte("download url=\"$(filter var=url)\"")
	output := string(expander{}.apply(in))
	t.Log(output)

	if string(output) != string(out) {
		t.Error("Out does not match.")
	}

	// In this case, no expansion should take place
	// download url="$(filter var=url)" file="0.jpg"
	in = []byte("download url=\"$(filter var=url)\"")
	out = []byte("download url=\"$(filter var=url)\"")

	if string(expander{}.apply(in)) != string(out) {
		t.Error("Out does not match.")
	}

	// download url="$(`url`)" file="0.jpg"
	in = []byte("download url=\"$(`url`)\" file=\"$(`i`).jpg\"")
	t.Log(string(expander{}.apply(in)))

	// download file="$(`0`).jpg"
	in = []byte("download file=\"$(`0`).jpg\"")
	t.Log(string(expander{}.apply(in)))

	// filter var=result query="entry_data.TagPage[0].tag.media.nodes[$(`i`)].display_src" => url
	in = []byte("filter var=result query=\"entry_data.TagPage[0].tag.media.nodes[$(`i`)].display_src\" => url")
	t.Log(string(expander{}.apply(in)))

	in = []byte("$(query query_id=\"17851374694183129\" variables={\"id\":\"$(`userId`)\",\"first\":$(`followersCount`)} => followers)")
	t.Log(string(expander{}.apply(in)))

}
