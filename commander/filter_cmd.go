package commander

import (
	"container/list"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/jmoiron/jsonq"
	"github.com/wsxiaoys/terminal/color"
)

// Filter ...
func (c *Commander) Filter(command string, tokens []string, data map[string]string) interface{} {
	filterVariable := strings.TrimSpace(data["var"])
	if len(data) == 0 || filterVariable == "" {
		color.Println("@r", command, "var=result query=entry_data.TagPage[0].tag.media.nodes[0].display_src")
		return nil
	}
	if len(filterVariable) == 1 && filterVariable[0] == '#' {
		return nil
	}
	isArray := filterVariable[0] == '#'
	var result interface{}
	if isArray {
		filterVariable = filterVariable[1:]
		if val, ok := c.Store[filterVariable].([]interface{}); ok {
			return len(val) - 1
		}
		return 0
	}
	result = c.Store[filterVariable]
	if result, ok := result.(string); ok {
		return result
	}
	if result, ok := result.(int); ok {
		return result
	}
	if result, ok := result.(bool); ok {
		return result
	}
	if result, ok := result.(map[string]interface{}); ok {
		return _jsonQuery(result, data)
	}
	if result, ok := result.([]interface{}); ok {
		index, err := strconv.Atoi(data["query"])
		if err != nil {
			return fmt.Sprintf("%v", result)
		}
		// Check bounds
		if index < 0 {
			index = 0
		}
		if index >= len(result) {
			index = len(result) - 1
		}
		return result[index]
	}
	return ""
}

func _jsonQuery(result interface{}, data map[string]string) interface{} {
	if result == nil {
		return map[string]interface{}{}
	}
	jq := jsonq.NewQuery(result)
	qComponents := func(in string) []string {
		return strings.Split(in, ".")
	}
	// Process array notation ex: TagPage[0], ProfilePage[0]
	qProcessArray := func(in []string) []string {
		// Regex to find array notation ex: TagPage[0], TagPage[1], TagPage[10]
		rgx, err := regexp.Compile(`\[.+\]`)
		if err != nil {
			return in
		}
		lst := list.New()
		for _, token := range in {
			if rgx.MatchString(token) {
				match := rgx.FindString(token)
				// Tag[0] -> Tag
				token = strings.Replace(token, match, "", -1)
				// Nex: [0] -> 0
				match = strings.Replace(match, "[", "", -1)
				match = strings.Replace(match, "]", "", -1)
				// Append match to our input array
				lst.PushBack(token)
				lst.PushBack(match)
			} else {
				lst.PushBack(token)
			}
		}
		out := []string{}
		for e := lst.Front(); e != nil; e = e.Next() {
			out = append(out, e.Value.(string))
		}
		return out
	}

	queryTokens := qProcessArray(qComponents(data["query"]))
	var obj interface{}
	var err error
	obj, err = jq.Object(queryTokens...)
	if err != nil {
		obj, err = jq.String(queryTokens...)
		if err != nil {
			obj, err = jq.Int(queryTokens...)
			if err != nil {
				obj, err = jq.Bool(queryTokens...)
			}
		}
	}
	if err != nil {
		log.Println(err)
		return nil
	}
	return obj
}
