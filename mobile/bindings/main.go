package bindings

import (
	"log"

	"github.com/hemantasapkota/insta/mobile/bindings/ios"
)

var appIos *IosApp

//SwitchState ...
func SwitchState(state string, options []byte) (bool, error) {
	log.SetPrefix("InstaMobile")
	log.Println("AppState changed to ", state)

	switch state {
	case ios.StateDidFinishLaunchingWithOptions:
		return appIos.appDidFinishLaunchingWithOptions(options)

	case ios.StateDidRegisterForRemoteNotificationsWithDeviceToken:
		appIos.appDidRegisterForRemoteNotificationsWithDeviceToken(nil)

	case ios.StateDidBecomeActive:
		appIos.appDidBecomeActive()

	case ios.StateDidEnterBackground:
		appIos.appDidEnterBackground()

	case ios.StateWillEnterForeground:
		appIos.appWillEnterForeground()

	case ios.StateAppOpenURL:
		return appIos.appOpenURL(options)

	case ios.StateWillTerminate:
		appIos.appWillTerminate()
	}

	return true, nil
}
