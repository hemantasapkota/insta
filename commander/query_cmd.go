package commander

import (
	"encoding/json"
        "log"
	"github.com/wsxiaoys/terminal/color"
)

// Filter ...
func (c *Commander) QueryCmd(command string, tokens []string, data map[string]string) interface{} {
	intent := c.Intents[command].(map[interface{}]interface{})
	if len(data) == 0 {
		color.Println("@r ", command, intent["Usage"])
		return nil
	}
	endpoint := c.makeEndpoint(intent["Endpoint"].(string), data)
        endpoint += "?" + formEncode(data)
	var body []byte
        _, body, _ = c.bot.Requester("GET", endpoint).Client.EndBytes()
	var m interface{}
	err := json.Unmarshal(body, &m)
	if err != nil {
		log.Println(err)
		return nil
	}
	c.Responses[command] = m
	return m
}
