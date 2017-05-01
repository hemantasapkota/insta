package commander

import (
	"fmt"
	"regexp"
	"strings"

	exp "instacli/expression"
)

func (c *Commander) unwrapQuotes(in string) (out string) {
	if in == "" {
		return
	}
	val := []byte(in)
	if len(val) > 0 {
		if val[0] == '"' {
			val[0] = ' '
		}
		if val[len(val)-1] == '"' {
			val[len(val)-1] = ' '
		}
		out = string(val)
	}
	out = strings.TrimSpace(out)
	return
}

func (c *Commander) parseCommandData(in string) (cmdStr string, dataString string) {
	// Find data part of the command. Data section starts with the first assignment (=)

	// Input command:
	// repeat frequency=10 cmd="$(like id=$(last_response cmd=scrape_entry_data query=entry_data.TagPage[0].media.nodes[$(COUNTER)]))"

	// Output:
	// frequency=10 cmd="like id=$(last_response cmd=scrape_entry_data query=entry_data.TagPage[0].media.nodes[$COUNTER])"

	// https://regex-golang.appspot.com/assets/html/index.html

	regex, _ := regexp.Compile(`\w+\=.+`)
	dataToken := regex.FindStringSubmatch(in)

	if len(dataToken) == 1 {
		dataString = dataToken[0]

		// replace the data part of the command
		cmdStr = strings.Replace(in, dataString, "", -1)

	} else {
		cmdStr = in
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
	// conv_id=jondoe_with_laex.pearl&text="Hello World"&"data="Hello"
	// [conv_id=jondoe_with_laex.pearl, text="Hello World", data="Hello"]

	for _, dataToken := range dataTokens {
		// We only want to split the string into two halves
		// Ex: value="$(last_response cmd=scrape_entry_data)"
		// => [value, "$(last_response cmd=scrape_entry_data)"]
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

func (c *Commander) parseCommand(in string) (string, []string, map[string]string) {
	// Commands can be nested.

	// $() => function
	// $COUNTER -> variable

	// Sample command:
	// repeat frequency=10 cmd="like id=$(last_response cmd=scrape_entry_data query=entry_data.TagPage[0].media.nodes[$(COUNTER)])"

	// Output: command: repeat
	//	   tokens: [repeat]
	//	   data: [cmd:like, frequency:10, id="$(last_response cmd=scrape_entry_data query=entry_data.TagPage[0].media.nodes[$(COUNTER)])"]

	// Regex tester
	// https://regex-golang.appspot.com/assets/html/index.html

	cmd, dataString := c.parseCommandData(in)
	data := c.processCommandData(dataString)

	// Evaluate nested commands
	newData := map[string]string{}
	for key, value := range data {
		if len(value) >= 3 && value[0] == '$' && value[1] == '(' {
			out := exp.Eval(exp.Parse(value), "", "", func(in string) string {
				return fmt.Sprintf("%v", c.Execute(in))
			})
			data[key] = out
		}
	}

	newData = data

	tokens := strings.Split(cmd, " ")
	return tokens[0], tokens, newData
}
