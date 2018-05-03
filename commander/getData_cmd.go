package commander

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/wsxiaoys/terminal/color"
)

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

