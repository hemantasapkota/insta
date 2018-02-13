package commander

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/wsxiaoys/terminal/color"
)

func formEncode(data map[string]string) string {
	values := make([]string, 0)
	for key, value := range data {
		values = append(values, fmt.Sprintf(`%s=%s`, key, value))
	}
	return strings.Join(values, "&")
}

// RequestExecutorCmd ...
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
	var m interface{}
	err := json.Unmarshal(body, &m)
	if err != nil {
		log.Println(err)
		return nil
	}
	c.Responses[command] = m
	return m
}

// GetDataCmd ...
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

// Counter ...
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

// Download ...
func (c *Commander) Download(command string, tokens []string, data map[string]string) interface{} {
	intent := c.Intents[command].(map[interface{}]interface{})
	if len(data) == 0 {
		color.Println("@r ", command, intent["Usage"])
		return nil
	}
	_url, ok := data["url"]
	if !ok || _url == "" {
		color.Println("@r ", command, intent["Usage"])
		return nil
	}
	_, body, errs := c.bot.Client.Get(_url).EndBytes()
	if errs != nil {
		color.Println("@r", command, errs[0])
		return nil
	}
	u, err := url.Parse(_url)
	if err != nil {
		color.Println("@r", command, errs[0])
		return nil
	}
	go func() {
		dir := filepath.Join(".", "downloads")
		_ = os.MkdirAll(dir, os.ModePerm)
		path := filepath.Join("./downloads", path.Base(u.Path))
		if path != "" {
			err := ioutil.WriteFile(path, body, 0644)
			if err != nil {
				// error writing the file
				color.Println("@r", command, err)
				return
			}
		}
	}()
	c.log.mediaMu.Lock()
	defer c.log.mediaMu.Unlock()
	c.log.Media[_url] = body
	c.log.Save(c.log)
	return "Downloaded " + _url
}
