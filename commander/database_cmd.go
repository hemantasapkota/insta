package commander

import (
	"github.com/wsxiaoys/terminal/color"
)

// SaveToDB ...
func (c *Commander) SaveToDB(command string, tokens []string, data map[string]string) (result interface{}) {
	if len(data) == 0 {
		color.Println("@r", command, "key= var= (Only variables can be stored)")
		return
	}
	key := data["key"]
	if key == "" {
		color.Println("@r", command, "key cannot be empty")
		return
	}
	c.log.logMu.Lock()
	defer c.log.logMu.Unlock()
	c.log.Log[key] = c.Store[data["val"]]
	c.log.Save(c.log)
	return nil
}

// ReadFromDB ...
func (c *Commander) ReadFromDB(command string, tokens []string, data map[string]string) (result interface{}) {
	if len(data) == 0 {
		color.Println("@r", command, "key=")
		return
	}
	return c.log.Log[data["key"]]
}
