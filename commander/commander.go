package commander

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/go-yaml/yaml"
	"github.com/peterh/liner"
	"github.com/wsxiaoys/terminal/color"
	"github.com/hemantasapkota/djangobot"
	"github.com/hemantasapkota/goma"
)

// A struct for logging the commands
var cmdLog = &commandLog{Log: map[string]interface{}{}, Media: map[string]interface{}{}}
var accountContext = ""

type commandLog struct {
	*goma.Object

	logMu sync.Mutex
	Log   map[string]interface{}

	mediaMu sync.Mutex
	Media   map[string]interface{}
}

func (log *commandLog) Key() string {
	if accountContext == "" {
		return "commandLog"
	}
	return accountContext
}

type cmdFunc func(command string, tokens []string, data map[string]string) interface{}

//Commander ...
type Commander struct {
	Intents   map[string]interface{}
	Responses map[string]interface{}
	Store     map[string]interface{}
	Commands  map[string]cmdFunc

	bot *djangobot.Bot
}

//New ...
func New(bot *djangobot.Bot) *Commander {
	commander := &Commander{
		Intents:   map[string]interface{}{},
		Responses: map[string]interface{}{},
		Store:     map[string]interface{}{},

		Commands: map[string]cmdFunc{},
		bot:      bot,
	}

	accountContext = bot.Username
	color.Print("@y> Press ? to list commands.")
	return commander
}

//LoadIntentsFromFile ...
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

//LoadIntents ...
func (c *Commander) LoadIntents(intents []byte) error {
	if intents == nil {
		return errors.New("intents data empty")
	}

	c.Intents = make(map[string]interface{})
	err := yaml.Unmarshal(intents, c.Intents)

	if err != nil {
		return err
	}

	for key := range c.Intents {
		c.Commands[key] = c.RequestExecutorCmd
	}

	// Load built in commands
	c.Commands["get_data"] = c.GetDataCmd
	c.Commands["filter"] = c.Filter
	c.Commands["run_script"] = c.RunScript
	c.Commands["counter"] = c.Counter
	c.Commands["download"] = c.Download

	// Experimental commands
	// c.Commands["repeat"] = c.Repeat

	return nil
}

//PrintCommands ...
func (c *Commander) PrintCommands() {
	for cmdName := range c.Commands {
		color.Print("@g\t ", cmdName)
		println()
	}
}

//Execute ...
func (c *Commander) Execute(command string) (result interface{}) {
	cmd, tokens, data, resultVar := c.parseCommand(command)

	// Execute the command
	functor, ok := c.Commands[cmd]
	if ok {
		result = functor(cmd, tokens, data)
		if result == nil {
			return
		}

		c.printYaml(result)

		// Store the result in c.Responses map
		c.Store[resultVar] = result

		// Log command and result
		intentObj, ok := c.Intents[cmd]
		if ok {
			intent := intentObj.(map[interface{}]interface{})
			if intent["Log"].(bool) {
				go func() {
					cmdLog.logMu.Lock()
					defer cmdLog.logMu.Unlock()

					cmdLog.Log[command] = result
					cmdLog.Save(cmdLog)
				}()
			}
		}
	}
	return
}

//Listen ...
func (c *Commander) Listen() {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	line.SetCompleter(func(line string) (list []string) {
		for cmdName := range c.Commands {
			if strings.HasPrefix(cmdName, line) {
				list = append(list, cmdName)
			}
		}
		return
	})

	historyFn := "liner_history"
	historyFile, err := os.Open(historyFn)
	if err != nil {
		historyFile, _ = os.Create(historyFn)
	}
	defer historyFile.Close()

	prompter := func() (string, error) {
		cmdName, err := line.Prompt("")
		cmdName = strings.TrimSpace(cmdName)
		if cmdName == "?" {
			c.PrintCommands()
		}
		return cmdName, err
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
