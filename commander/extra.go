package commander

import (
	"fmt"
	"github.com/wsxiaoys/terminal/color"
	"strconv"
	"time"
)

// Experimental
func (c *Commander) Repeat(command string, tokens []string, data map[string]string) (result interface{}) {
	if len(data) == 0 {
		color.Println("@r", command, " frequency=(in seconds) cmd=")
		return
	}

	frequency, err := strconv.Atoi(data["frequency"])
	if err != nil {
		fmt.Println("Frequency not set.")
		return
	}

	cmdToExec := data["cmd"]
	go func(in string) {
		remaining := frequency
		for {
			select {
			case <-time.After(time.Second * time.Duration(frequency)):
				if remaining == 0 {
					fmt.Println("Finished repeat.")
					return
				}
				c.Execute(in)
				remaining--
			}
		}
	}(cmdToExec)

	return
}
