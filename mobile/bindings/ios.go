package bindings

import (
	"encoding/json"
	"log"

	"github.com/hemantasapkota/goma/gomadb"
	ldb "github.com/hemantasapkota/goma/gomadb/leveldb"

	"github.com/hemantasapkota/insta/mobile/bindings/ios"
)

//IosApp ...
type IosApp struct {
	ios.UIApplicationDelegate
}

func (app *IosApp) appDidFinishLaunchingWithOptions(options []byte) (bool, error) {
	log.SetPrefix("InstaMobile: ")
	log.Println("Welcome to Go layer.")

	appOptions := map[string]interface{}{}
	err := json.Unmarshal(options, appOptions)
	if err != nil {
		return false, err
	}

	// Init our database
	db, err := ldb.InitDB(appOptions["dbPath"].(string))
	if err != nil {
		return false, err
	}
	gomadb.SetLevelDB(db)

	err = initApp()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (app *IosApp) appDidRegisterForRemoteNotificationsWithDeviceToken(deviceToken []byte) {
}

func (app *IosApp) appDidReceiveRemoteNotification(userInfo []byte) {
}

func (app *IosApp) appDidFailToRegisterForRemoteNotificationsWithError(error []byte) {
}

func (app *IosApp) appWillResignActive() {
}

func (app *IosApp) appDidEnterBackground() {
}

func (app *IosApp) appWillEnterForeground() {
}

func (app *IosApp) appDidBecomeActive() {
}

func (app *IosApp) appWillTerminate() {
}

func (app *IosApp) appOpenURL(options []byte) (bool, error) {
	return true, nil
}
