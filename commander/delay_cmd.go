package commander

import (
	"fmt"
	"strconv"

	"github.com/wsxiaoys/terminal/color"
)

// Delay ...
func (c *Commander) Delay(command string, tokens []string, data map[string]string) (result interface{}) {
	if len(data) == 0 {
		color.Println("@r", command, " interval=(in seconds)")
		return
	}
	delay, err := strconv.Atoi(data["interval"])
	if err != nil {
		return nil
	}
	c.cmdDelay = delay
	return fmt.Sprintf("Delay set to %d", c.cmdDelay)
}
