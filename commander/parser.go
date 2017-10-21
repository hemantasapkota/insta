package commander

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	exp "github.com/hemantasapkota/insta/expression"
)

func toInt(index string) int {
	i, err := strconv.Atoi(index)
	if err != nil {
		return 0
	}
	return i
}

func (c *Commander) unwrapQuotes(in string) (out string) {
	if in == "" {
		return
	}
	val := []byte(in)
	if len(val) > 0 {
		if val[0] == '"' && val[len(val)-1] == '"' {
			val[0], val[len(val)-1] = ' ', ' '
		}
		out = string(val)
	}
	out = strings.TrimSpace(out)
	return
}

func (c *Commander) parseCommandData(in string) (cmdStr string, dataString string, resultStr string) {
	// Input command:
	// repeat frequency=10 cmd="$(like id=$(last_response cmd=scrape_entry_data query=entry_data.TagPage[0].media.nodes[$(COUNTER)]))"
	// Output:
	// frequency=10 cmd="like id=$(last_response cmd=scrape_entry_data query=entry_data.TagPage[0].media.nodes[$COUNTER])"

	// https://regex-golang.appspot.com/assets/html/index.html
	leftRight := strings.Split(in, "=>")
	if len(leftRight) > 1 {
		resultStr = strings.TrimSpace(leftRight[1])
	}

	regex, _ := regexp.Compile(`\w+\=.+`)

	cmd := leftRight[0]
	dataToken := regex.FindStringSubmatch(cmd)

	if len(dataToken) == 1 {
		dataString = dataToken[0]
		// replace the data part of the command
		cmdStr = strings.Replace(cmd, dataString, "", -1)
	} else {
		cmdStr = cmd
	}

	return
}

func (c *Commander) processCommandData(in string) map[string]string {
	data := map[string]string{}

	// Replace all spaces with &, but don't replace spaces inside " "
	// This: conv_id=jondoe_with_laex.pearl text="Hello World." data="Hello"
	// Becomes: conv_id=jondoe_with_laex.pearl&text="Hello World"&data="Hello"
	token := []byte(in)
	counter := 0
	i := 0
	for i < len(token) {
		c := token[i]
		if c == ' ' {
			if counter > 0 {
				i++
				continue
			} else {
				token[i] = byte('&')
			}
		}
		if string(c) == `"` {
			if counter > 0 {
				counter = 0
			} else {
				counter = 1
			}
		}
		i++
	}

	dataTokens := strings.Split(string(token), "&")
	for _, dataToken := range dataTokens {
		// We only want to split the string into two halves
		// Ex: value="$(last_response cmd=scrape_entry_data)" => [value, "$(last_response cmd=scrape_entry_data)"]
		// In the SplitN, the parameter 2 means only return two substrings
		items := strings.SplitN(dataToken, "=", 2)

		if len(items) == 2 {
			key := items[0]
			// The value may be wrapped in quotes. So we replace them
			data[key] = c.unwrapQuotes(items[1])
		}
	}

	return data
}

func (c *Commander) parseCommand(in string) (string, []string, map[string]string, string) {
	cmd, dataString, resultVar := c.parseCommandData(in)
	data := c.processCommandData(dataString)
	tokens := strings.Split(cmd, " ")

	node := exp.Parse(in).Prune()
	expChecker := &exp.ExpChecker{Node: node}

	if expChecker.IsLoop() {
		exp.Eval(node, "", "", map[string]string{}, func(inexp string) string {
			return fmt.Sprintf("%v", c.Execute(inexp))
		})
		return tokens[0], tokens, data, resultVar
	}

	if expChecker.IsPool() {
		c.loop.process(c)
		exp.Eval(node, "", "", map[string]string{}, func(inexp string) string {
			return fmt.Sprintf("%v", c.Execute(inexp))
		})
		return tokens[0], tokens, data, resultVar
	}

	if c.loop != nil {
		c.loop.batch = append(c.loop.batch, in)
		return tokens[0], tokens, data, resultVar
	}

	// Non looping execution
	exp.Eval(node, "", "", map[string]string{}, func(inexp string) string {
		return fmt.Sprintf("%v", c.Execute(inexp))
	})

	return tokens[0], tokens, data, resultVar
}
