package commander

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jmoiron/jsonq"
	"github.com/wsxiaoys/terminal/color"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func formEncode(data map[string]string) string {
	values := make([]string, 0)
	for key, value := range data {
		values = append(values, fmt.Sprintf(`%s=%s`, key, value))
	}
	return strings.Join(values, "&")
}


func (c *Commander) RequestExecutorCmd(command string, tokens []string, data map[string]string) interface{} {
	intent := c.Intents[command].(map[interface{}]interface{})
	if len(data) == 0 {
		color.Println("@r ", command, intent["Usage"])
		return nil
	}

	// Make template
	endpoint := c.makeEndpoint(intent["Endpoint"].(string), data)

	var body []byte

	// get or post
	method := intent["Method"].(string)

	switch method {
	case "get":
		_, body, _ = c.bot.Requester("GET", endpoint).Client.EndBytes()
	case "post":
		_, body, _ = c.bot.Requester("POST", endpoint).Client.Send(formEncode(data)).EndBytes()
	case "custom":
		// For custom methods we forward execution
		cmd, ok := c.Commands[command]
		if ok {
			cmd(command, tokens, data)
		}
	}

	// Unmarshal response
	var m interface{}
	err := json.Unmarshal(body, &m)

	if err != nil {
		log.Println(err)
		return nil
	}

	c.Responses[command] = m

	return m
}

func (c *Commander) ScrapeEntryDataCmd(command string, tokens []string, data map[string]string) interface{} {
	intent := c.Intents[command].(map[interface{}]interface{})
	if len(data) == 0 {
		color.Println("@r ", command, intent["Usage"])
		return nil
	}

	endpoint := c.makeEndpoint(intent["Endpoint"].(string), data)

	// Scrape the url
	doc, err := goquery.NewDocument(endpoint)
	if err != nil {
		log.Println(err)
		return nil
	}

	var m interface{}

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if !strings.Contains(text, "window._sharedData") {
			return
		}

		text = strings.Replace(text, "window._sharedData = ", "", -1)
		text = text[0 : len(text)-1]

		err = json.Unmarshal([]byte(text), &m)
		if err != nil {
			fmt.Println("Failed to scrape entry data.")
			return
		}

		c.Responses[command] = m
	})

	return m
}

func (c *Commander) LastResponseCmd(command string, tokens []string, data map[string]string) interface{} {
	if len(data) == 0 {
		color.Println("@r", command, "cmd=scrape_entry_data query=entry_data.TagPage[0].tag.media.page_info.end_cursor")
		return nil
	}

	response, ok := c.Responses[data["cmd"]]
	if !ok {
		return nil
	}

	// JSON Object -> Query
	jq := jsonq.NewQuery(response)

	// Get our query components
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

		// Go through the list
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
	}

	if err != nil {
		log.Println(err)
		return nil
	}

	return obj
}

func (c *Commander) Counter(command string, tokens []string, data map[string]string) interface{} {
	updateVal, ok := data["set"]
	if ok {
		intVal, err := strconv.Atoi(updateVal)
		if err != nil {
			intVal = 0
		}
		c.Store[command] = intVal
		return updateVal
	}

	val, ok := c.Store[command]
	if ok {
		val = val.(int) + 1
		c.Store[command] = val
		return val
	}

	c.Store[command] = 0
	return 0
}


