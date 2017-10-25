package commander

import (
	"strconv"
	"strings"

	"github.com/wsxiaoys/terminal/color"
)

// Loop ...
func (c *Commander) Loop(command string, tokens []string, data map[string]string) (result interface{}) {
	if len(data) == 0 {
		color.Println("@r", command, "range=1,10 => i")
		return
	}
	loopRange := data["range"]
	items := strings.Split(loopRange, ",")
	if len(items) != 2 {
		color.Println("@r", "range improper.")
		return
	}
	start, _ := strconv.Atoi(items[0])
	end, _ := strconv.Atoi(items[1])
	if start > end {
		color.Println("@r", "only forward loops allowed.")
		return
	}
	result = start
	// setup loop context
	c.loop = &loopCtx{
		index:    start,
		endIndex: end,
		varName:  "i",
	}
	return
}

// Pool ...
func (c *Commander) Pool(command string, tokens []string, data map[string]string) (result interface{}) {
	c.loop = nil
	return ""
}
