package bindings

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/hemantasapkota/djangobot"
	"github.com/hemantasapkota/goma/gomadb"
	"github.com/hemantasapkota/insta/commander"
)

var cmdHandler *commander.Commander

func initApp() error {

	db := gomadb.GetDB()

	user, err := db.Get("username")
	if err != nil {
		return errors.New("username not found")
	}

	password, err := db.Get("password")
	if err != nil {
		return errors.New("password not found")
	}

	instabot := djangobot.With("https://www.instagram.com/accounts/login/ajax/").
		ForHost("instagram.com").
		SetUsername(user).
		SetPassword(password).
		LoadCookies()

	if instabot.Error != nil {
		return instabot.Error
	}

	_, err = instabot.
		X("csrfmiddlewaretoken", instabot.Cookie("csrftoken").Value).
		X("username", instabot.Username).
		X("password", instabot.Password).Login()

	sessionid := instabot.Cookie("sessionid").Value
	if sessionid == "" {
		return errors.New("Authentication failed with Instagram")
	}

	// Setup our commander
	intents, err := db.GetBytes("intents")
	if err != nil {
		return errors.New("intents data not found")
	}

	cmdHandler = commander.New(instabot)
	cmdHandler.LoadIntents([]byte(intents))

	return nil
}

//Execute ...
func Execute(cmd string) ([]byte, error) {
	log.Println("Executing ", cmd)
	result := cmdHandler.Execute(cmd)
	data, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return data, nil
}
