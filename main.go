package main

import (
	"flag"
	"fmt"
	"github.com/wsxiaoys/terminal/color"
	"github.com/hemantasapkota/goma/gomadb"
	ldb "github.com/hemantasapkota/goma/gomadb/leveldb"
	"github.com/hemantasapkota/djangobot"
	"instacli/commander"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"strings"
)

const usage = `
Usage:
	-username -password
	-account ( If specified in the .credentials.yaml file )
`

func main() {

	username := flag.String("username", "", "Username")
	password := flag.String("password", "", "Password")
	account := flag.String("account", "", "Account from .credentials.yaml")

	flag.Parse()

	if *username == "" || *password == "" {
		data, err := ioutil.ReadFile(".credentials.yaml")
		if err != nil {
			fmt.Println(usage)
			return
		}

		credentials := make(map[string]interface{})
		err = yaml.Unmarshal(data, credentials)

		if err != nil {
			fmt.Println(usage)
			return
		}

		if len(credentials) == 0 {
			fmt.Println(usage)
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
			for key, _ := range credentials {
				accountName = key
				user, pass := getAccountCreds(key)
				username = &user; password = &pass;
			}
			color.Println("@yAuthenticating with account: ", strings.TrimSpace(accountName))
		} else {
			if *account == "" {
				fmt.Println(usage)
				return
			}

			user, pass := getAccountCreds(*account)
			username = &user; password = &pass;
		}
	}

	// Instagram
	instabot := djangobot.With("https://www.instagram.com/accounts/login/ajax/").
			ForHost("instagram.com").
			SetUsername(*username).
			SetPassword(*password).
			LoadCookies()

	if instabot.Error != nil {
		panic(instabot.Error)
	}

	_, err := instabot.
		X("csrfmiddlewaretoken", instabot.Cookie("csrftoken").Value).
		X("username", instabot.Username).
		X("password", instabot.Password).Login()

	sessionid := instabot.Cookie("sessionid").Value
	if sessionid == "" {
		color.Println("@r Authentication failed with Instagram.")
		return
	}

	// Init our database
	db, err := ldb.InitDB(".")
	if err != nil {
		panic("Couldn't init database")
	}
	gomadb.SetLevelDB(db)

	// Setup our commander
	commandHandler := commander.New(instabot).LoadIntentsFromFile("instagram.yaml")

	commandHandler.Listen()
}
