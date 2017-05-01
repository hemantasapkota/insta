package commander

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/peterh/liner"
	"github.com/wsxiaoys/terminal/color"
	"goma"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"github.com/hemantasapkota/djangobot"
)

// A struct for logging the commands
var commandLog CommandLog = CommandLog{Log: map[string]interface{}{}}
var accountContext string = ""

type CommandLog struct {
	*goma.Object
	sync.Mutex

	Log map[string]interface{}
}

func (log CommandLog) Key() string {
	if accountContext == "" {
		return "commandLog"
	}
	return accountContext
}

type CmdFunc func(command string, tokens []string, data map[string]string) interface{}

type Commander struct {
	Intents   map[string]interface{}
	Responses map[string]interface{}
	Store     map[string]interface{}
	Commands  map[string]CmdFunc

	bot *djangobot.Bot
}

func New(bot *djangobot.Bot) *Commander {
	commander := &Commander{
		Intents:   map[string]interface{}{},
		Responses: map[string]interface{}{},
		Store:     map[string]interface{}{},

		Commands: map[string]CmdFunc{},
		bot:      bot,
	}

	accountContext = bot.Username
	color.Print("@y> Press ? to list commands.")
	return commander
}

func (c *Commander) LoadIntentsFromFile(filename string) *Commander {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return c
	}
	c.LoadIntents(data)
	return c
}

func (c *Commander) printYaml(m interface{}) {
	yamlBody, err := yaml.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s", string(yamlBody))
}

func (c *Commander) makeTemplate(tmpl string, params interface{}) (string, error) {
	t := template.Must(template.New("template").Parse(tmpl))
	buf := &bytes.Buffer{}
	err := t.Execute(buf, params)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *Commander) makeEndpoint(endpoint string, data map[string]string) string {
	endpoint, _ = c.makeTemplate(endpoint, struct {
		ID       string
		Query    string
		Path     string
		UserName string
	}{
		data["id"],
		data["query"],
		data["path"],
		c.bot.Username,
	})
	return endpoint
}

func (c *Commander) LoadIntents(intents []byte) error {
	if intents == nil {
		return errors.New("Intents data empty.")
	}

	c.Intents = make(map[string]interface{})
	err := yaml.Unmarshal(intents, c.Intents)

	if err != nil {
		return err
	}

	for key, _ := range c.Intents {
		c.Commands[key] = c.RequestExecutorCmd
	}

	// Load built in commands
	c.Commands["last_response"] = c.LastResponseCmd
	c.Commands["scrape_entry_data"] = c.ScrapeEntryDataCmd
	c.Commands["counter"] = c.Counter

	// Experimental commands
	//c.Commands["repeat"] = c.Repeat
	//c.Commands["map"] = c.Map

	return nil
}

func (c *Commander) PrintCommands() {
	for cmdName, _ := range c.Commands {
		color.Print("@g\t ", cmdName)
		println()
	}
}

func (c *Commander) Execute(command string) (result interface{}) {
	cmd, tokens, data := c.parseCommand(command)

	// Execute the command
	functor, ok := c.Commands[cmd]
	if ok {
		result = functor(cmd, tokens, data)
		if result == nil {
			return
		}

		c.printYaml(result)

		// Log command and result
		intentObj, ok := c.Intents[cmd]
		if ok {
			intent := intentObj.(map[interface{}]interface{})
			if intent["Log"].(bool) {
				go func() {
					commandLog.Lock()
					defer commandLog.Unlock()

					commandLog.Log[command] = result
					commandLog.Save(commandLog)
				}()
			}
		}
	}

	return
}

func (c *Commander) Listen() {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	line.SetCompleter(func(line string) (list []string) {
		for cmd_name, _ := range c.Commands {
			if strings.HasPrefix(cmd_name, line) {
				list = append(list, cmd_name)
			}
		}
		return
	})

	history_fn := "liner_history"
	history_file, err := os.Open(history_fn)
	if err != nil {
		history_file, _ = os.Create(history_fn)
	}
	defer history_file.Close()

	prompter := func() (string, error) {
		cmd_name, err := line.Prompt("")
		cmd_name = strings.TrimSpace(cmd_name)
		if cmd_name == "?" {
			c.PrintCommands()
		}
		return cmd_name, err
	}

	for {
		command, err := prompter()
		if err == liner.ErrPromptAborted {
			panic("Exiting")
		} else {
			go c.Execute(command)
			line.AppendHistory(command)
			color.Print("@y>")
		}
	}

}
