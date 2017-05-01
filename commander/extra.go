package commander

import (
	"fmt"
	"time"
	"github.com/wsxiaoys/terminal/color"
	"strconv"
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

func (c *Commander) Map(command string, tokens []string, data map[string]string) interface{} {
	if len(data) == 0 || data["key"] == "" {
		color.Println("@r", command, "key= value=")
		return nil
	}

	if data["value"] == "" {
		val := c.Store[data["key"]]
		c.printYaml(val)
		return val
	}

	// Store value
	value := data["value"]
	c.Store[data["key"]] = value
	return value
}
