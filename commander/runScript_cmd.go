package commander

import (
	"encoding/base64"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/wsxiaoys/terminal/color"
)

// RunScript ...
func (c *Commander) RunScript(command string, token []string, data map[string]string) interface{} {
	if len(data) == 0 {
		color.Println("@r", command, "file= fromLine=10 ( Ex: Start execution from line 10 )")
		return nil
	}
	file, ok := data["file"]
	lineFrom, fromOk := data["lineFrom"]
	lineTo, toOk := data["lineTo"]
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
		var start, end = 0, len(scripts)
		if fromOk {
			val, err := strconv.Atoi(lineFrom)
			if err == nil {
				// ex: for line 20, start from 19
				start = val - 1
				if start > len(scripts) {
					start = len(scripts) - 1
				}
				if start < 0 {
					start = 0
				}
			}
		}
		if toOk {
			val, err := strconv.Atoi(lineTo)
			if err != nil {
				end = len(scripts) - 1
			} else {
				end = val
			}
		}
		// Execute
		for i := start; i < end; i++ {
			statement := strings.TrimSpace(scripts[i])
			if statement != "" && statement[0] != '#' {
				c.Execute(statement)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
	return nil
}

// RunScriptBase64 ...
func (c *Commander) RunScriptBase64(command string, token []string, data map[string]string) interface{} {
	if len(data) == 0 {
		color.Println("@r", command, `script="JChsb29wIHJhbmdlPTEsMTAgPT4gaSkKICAkKCQoYGlgKSA9PT4gbGlzdCkKJChwb29sKQo=" (data encoded in base64)`)
		return nil
	}
	script := data["script"]
	decoded, err := base64.StdEncoding.DecodeString(script)
	if err != nil {
		color.Println("@r", command, "Invalid base64 data.")
	}
	dataStr := strings.TrimSpace(string(decoded))
	scripts := strings.Split(dataStr, "\n")
	var i = 0
	for ; i < len(scripts); i++ {
		statement := strings.TrimSpace(scripts[i])
		if statement != "" && statement[0] != '#' {
			c.Execute(statement)
			time.Sleep(100 * time.Millisecond)
		}
	}
	return nil
}
