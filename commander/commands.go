package commander

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jmoiron/jsonq"
	"github.com/wsxiaoys/terminal/color"
)

func formEncode(data map[string]string) string {
	values := make([]string, 0)
	for key, value := range data {
		values = append(values, fmt.Sprintf(`%s=%s`, key, value))
	}
	return strings.Join(values, "&")
}

//RequestExecutorCmd ...
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

//GetDataCmd ...
func (c *Commander) GetDataCmd(command string, tokens []string, data map[string]string) interface{} {
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

func _jsonQuery(result interface{}, data map[string]string) interface{} {
	if result == nil {
		return map[string]interface{}{}
	}

	// JSON Object -> Query
	jq := jsonq.NewQuery(result)

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

func (c *Commander) Filter(command string, tokens []string, data map[string]string) interface{} {
	if len(data) == 0 || data["var"] == "" {
		color.Println("@r", command, "var=result query=entry_data.TagPage[0].tag.media.nodes[0].display_src")
		return nil
	}

	result := c.Store[data["var"]]

	if result, ok := result.(string); ok {
		return result
	}

	return _jsonQuery(result, data)
}

func (c *Commander) RunScript(command string, token []string, data map[string]string) interface{} {
	if len(data) == 0 {
		color.Println("@r", command, "file=")
		return nil
	}

	file, ok := data["file"]
	if ok {
		// path =
		script := filepath.Join(".", file)
		data, err := ioutil.ReadFile(script)
		if err != nil {
			color.Println("@r", command, err)
			return nil
		}

		dataStr := strings.TrimSpace(string(data))
		if len(dataStr) == 0 {
			color.Println("@r", command, "Script Empty.")
			return nil
		}

		scripts := strings.Split(dataStr, "\n")
		for _, statement := range scripts {
			// ignore empty strings or comments
			stmt := strings.TrimSpace(statement)
			if stmt != "" || stmt[0] != '#' {
				c.Execute(statement)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	return nil
}

//Counter ...
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

//Download ...
func (c *Commander) Download(command string, tokens []string, data map[string]string) interface{} {
	intent := c.Intents[command].(map[interface{}]interface{})
	if len(data) == 0 {
		color.Println("@r ", command, intent["Usage"])
		return nil
	}

	url, ok := data["url"]
	if !ok {
		color.Println("@r ", command, intent["Usage"])
		return nil
	}

	_, body, errs := c.bot.Client.Get(url).EndBytes()
	if errs != nil {
		color.Println("@r", command, errs[0])
		return nil
	}

	file, ok := data["file"]
	if ok {
		go func() {
			dir := filepath.Join(".", "downloads")
			_ = os.MkdirAll(dir, os.ModePerm)
			path := filepath.Join("./downloads", file)
			if path != "" {
				err := ioutil.WriteFile(path, body, 0644)
				if err != nil {
					// error writing the file
					color.Println("@r", command, err)
					return
				}
			}
		}()
	}

	cmdLog.mediaMu.Lock()
	defer cmdLog.mediaMu.Unlock()

	cmdLog.Media[url] = body
	cmdLog.Save(cmdLog)

	return "Success."
}
