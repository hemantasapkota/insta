package commander

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/hemantasapkota/djangobot"
	"github.com/hemantasapkota/goma"
	exp "github.com/hemantasapkota/insta/expression"
	"github.com/hemantasapkota/insta/flags"
	"github.com/peterh/liner"
	"github.com/wsxiaoys/terminal/color"
)

// A struct for logging the commands
var accountContext = ""

// loop context
type loopCtx struct {
	index    int
	endIndex int
	varName  string
	batch    []string
}

func (ctx *loopCtx) process(commander *Commander) {
	for ctx.index <= ctx.endIndex {
		for _, stmt := range ctx.batch {
			node := exp.Parse(stmt).Prune()
			exp.Eval(node, "", "", map[string]string{}, func(inexp string) string {
				return fmt.Sprintf("%v", commander.Execute(inexp))
			})
		}
		nextIndex := ctx.index + 1
		commander.Store[ctx.varName] = nextIndex
		ctx.index++
	}
}

// command log
type commandLog struct {
	*goma.Object
	logMu   sync.Mutex
	Log     map[string]interface{}
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

// Commander ...
type Commander struct {
	Intents   map[string]interface{}
	Responses map[string]interface{}
	Store     map[string]interface{}
	Commands  map[string]cmdFunc

	log      *commandLog
	loop     *loopCtx
	bot      *djangobot.Bot
	cmdDelay int
}

// New ...
func New(bot *djangobot.Bot) *Commander {
	commander := &Commander{
		Intents:   map[string]interface{}{},
		Responses: map[string]interface{}{},
		Store:     map[string]interface{}{},
		Commands:  map[string]cmdFunc{},
		bot:       bot,
		log:       &commandLog{Log: map[string]interface{}{}, Media: map[string]interface{}{}},
	}
	accountContext = bot.Username
	if !flags.Silent {
		color.Print("@y> Press ? to list commands.")
	}
	return commander
}

// LoadIntentsFromFile ...
func (c *Commander) LoadIntentsFromFile(filename string) *Commander {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return c
	}
	c.LoadIntents(data)
	return c
}

func (c *Commander) printOutput(m interface{}) {
	if flags.OutputFormat == "yaml" {
		yamlBody, err := yaml.Marshal(m)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s", string(yamlBody))
		return
	}
	if flags.OutputFormat == "json" {
		jsonBody, err := json.MarshalIndent(m, "", "    ")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s", string(jsonBody))
	}
	if flags.ExecFile {
		println()
	}
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

// LoadIntents ...
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
	c.Commands["query"] = c.QueryCmd
	c.Commands["filter"] = c.Filter
	c.Commands["run_script"] = c.RunScript
	c.Commands["run_script_base64"] = c.RunScriptBase64
	c.Commands["counter"] = c.Counter
	c.Commands["download"] = c.Download
	c.Commands["save_to_db"] = c.SaveToDB
	c.Commands["loop"] = c.Loop
	c.Commands["pool"] = c.Pool
	c.Commands["delay"] = c.Delay
	return nil
}

// PrintCommands ...
func (c *Commander) PrintCommands() {
	for cmdName := range c.Commands {
		color.Print("@g\t ", cmdName)
		println()
	}
}

func evalIfBlock(ifBlock string) (bool, error) {
	components := strings.Split(ifBlock, "_")
	if !(len(components) == 3) {
		return false, errors.New("if block should have three components")
	}
	var result bool
	switch components[1] {
	case "contains":
		result = strings.Contains(strings.ToLower(components[0]), strings.ToLower(components[2]))
	case "equals":
		result = strings.ToLower(components[0]) == strings.ToLower(components[2])
	case "notEquals":
		result = strings.ToLower(components[0]) != strings.ToLower(components[2])
	case "lessThan":
		num1, _ := strconv.Atoi(strings.ToLower(components[0]))
		num2, _ := strconv.Atoi(strings.ToLower(components[2]))
		result = num1 < num2
	case "aboveThan":
		num1, _ := strconv.Atoi(strings.ToLower(components[0]))
		num2, _ := strconv.Atoi(strings.ToLower(components[2]))
		result = num1 > num2
	default:
		return false, errors.New("unknown if condition")
	}
	return result, nil
}

// Execute ...
func (c *Commander) Execute(command string) (result interface{}) {
	cmd, tokens, data, resultVar, assignType := c.parseCommand(strings.TrimSpace(command))
	functor, ok := c.Commands[cmd]
	if ok {
		if c.cmdDelay > 0 {
			time.Sleep(time.Second * time.Duration(c.cmdDelay))
			c.cmdDelay = 0
		}
		ifBlock, ok := data["if"]
		if ok {
			ifResult, err := evalIfBlock(ifBlock)
			if err != nil {
				return
			}
			if !ifResult {
				return
			}
		}
		result = functor(cmd, tokens, data)
		if result == nil {
			return
		}
		c.printOutput(result)
		if resultVar != "" {
			if assignType == "=>" {
				c.Store[resultVar] = result
			}
			// append
			if assignType == "==>" {
				_, ok := c.Store[resultVar]
				if !ok {
					c.Store[resultVar] = make([]interface{}, 0)
				}
				list := c.Store[resultVar].([]interface{})
				list = append(list, result)
				c.Store[resultVar] = list
			}
		}

		intentObj, ok := c.Intents[cmd]
		if ok {
			intent := intentObj.(map[interface{}]interface{})
			if intent["Log"].(bool) {
				go func() {
					c.log.logMu.Lock()
					defer c.log.logMu.Unlock()

					c.log.Log[command] = result
					c.log.Save(c.log)
				}()
			}
		}
	}
	return
}

// Listen ...
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
