package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

        "github.com/parnurzeal/gorequest"
	"github.com/hemantasapkota/djangobot"
	"github.com/hemantasapkota/goma/gomadb"
	ldb "github.com/hemantasapkota/goma/gomadb/leveldb"

	"github.com/hemantasapkota/insta/commander"
	"github.com/hemantasapkota/insta/flags"

	"github.com/wsxiaoys/terminal/color"
	"gopkg.in/yaml.v2"
)

func main() {
	username := flag.String("username", "", "Use \"test\" for testing.")
	password := flag.String("password", "", "Use \"test\" for testing.")
	account := flag.String("account", "", "Account from .credentials.yaml.")
	exec := flag.String("exec", "", "Execute a command")
	execFile := flag.String("execFile", "", "Execute file")
	silent := flag.Bool("silent", false, "Only ouputs and errors will be printed.")
	json := flag.Bool("json", false, "Output json")

	flag.Parse()

	flags.Silent = *silent
	if *json {
		flags.OutputFormat = "json"
	}
	if *username == "" || *password == "" {
		data, err := ioutil.ReadFile(".credentials.yaml")
		if err != nil {
			flag.PrintDefaults()
			return
		}
		credentials := make(map[string]interface{})
		err = yaml.Unmarshal(data, credentials)
		if err != nil {
			flag.PrintDefaults()
			return
		}
		if len(credentials) == 0 {
			flag.PrintDefaults()
			return
		}
		getAccountCreds := func(name string) (username string, password string) {
			account, ok := credentials[name].(map[interface{}]interface{})
			if ok {
				username = account["username"].(string)
				password = account["password"].(string)
			}
			return
		}
		if len(credentials) == 1 {
			accountName := ""
			for key := range credentials {
				accountName = key
				user, pass := getAccountCreds(key)
				username = &user
				password = &pass
			}
			if !flags.Silent {
				color.Println("@yAuthenticating with account: ", strings.TrimSpace(accountName))
			}
		} else {
			if *account == "" {
				flag.PrintDefaults()
				return
			}
			user, pass := getAccountCreds(*account)
			username = &user
			password = &pass
		}
	}

	// Test username
	if *username == "test" && *password == "test" {
		flags.IsTestnet = true
	}

	// Instagram
	var instabot *djangobot.Bot
	if flags.IsTestnet {
		instabot = &djangobot.Bot{
			Username: "test",
			Password: "test",
                        Client: gorequest.New(),
		}
	} else {
		instabot = djangobot.With("https://www.instagram.com/accounts/login/ajax/").
			ForHost("instagram.com").
			SetUsername(*username).
			SetPassword(*password).
			LoadCookies()
		if instabot.Error != nil {
			panic(instabot.Error)
		}
		instabot.
			X("csrfmiddlewaretoken", instabot.Cookie("csrftoken").Value).
			X("username", instabot.Username).
			X("password", instabot.Password).Login()
		sessionid := instabot.Cookie("sessionid").Value
		if sessionid == "" {
			color.Println("@r Authentication failed with Instagram.")
			return
		}
	}
	// Init our database
	db, err := ldb.InitDB(".")
	if err != nil {
		panic("Couldn't init database")
	}
	gomadb.SetLevelDB(db)
	// Setup our commander
	commandHandler := commander.New(instabot).LoadIntentsFromFile("instagram.yaml")
	// If exec mode, then execute a command and exit
	if *exec != "" {
		commandHandler.Execute(*exec)
		return
	}
	// execute script file
	if *execFile != "" {
		flags.ExecFile = true
		commandHandler.Execute(fmt.Sprintf(`run_script file="%s"`, *execFile))
		return
	}
	commandHandler.Listen()
}
